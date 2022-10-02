package processors

import (
	"fmt"
	"github.com/s4kibs4mi/jally-commerce-bot/config"
	"github.com/s4kibs4mi/jally-commerce-bot/log"
	"github.com/s4kibs4mi/jally-commerce-bot/models"
	"github.com/s4kibs4mi/jally-commerce-bot/models/api_request"
	"github.com/s4kibs4mi/jally-commerce-bot/services"
	"github.com/s4kibs4mi/jally-commerce-bot/services/messenger"
	"github.com/s4kibs4mi/jally-commerce-bot/services/messenger/template"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Action struct {
	Name         string
	Fn           func(senderID, template string, args []string)
	MenuPattern  string
	MenuTemplate string
	IsMenu       bool
}

type FacebookStateProcessor struct {
	stateService    services.IStateService
	shopemaaService services.IShopemaaService
	cfg             *config.Application
	messenger       *messenger.Messenger
}

func (p *FacebookStateProcessor) Init() error {
	if err := p.messenger.SetGetStartedButton("Say, Hello"); err != nil {
		log.Log().Errorln(err)
		os.Exit(-1)
	}
	if err := p.messenger.SetPersistentMenu(p.loadMenus()); err != nil {
		return err
	}
	return nil
}

func (p *FacebookStateProcessor) Process(req *api_request.CustomerRequest) error {
	log.Log().Infoln("Message Received....")
	if req.IsMessage {
		return p.processMessage(req.Event, req.Opts, req.Message)
	}
	return p.processPostback(req.Event, req.Opts, req.Postback)
}

func (p *FacebookStateProcessor) ProcessOrderCreated(cartID, orderHash, email string) error {
	to, _ := p.stateService.GetIdentityByCartID("to", cartID)

	msg := fmt.Sprintf("Thank you for the order\n. OrderHash: %s\n. We will process it accordingly.", orderHash)
	orderDetailsBtn := template.NewWebURLButton("Order Details", fmt.Sprintf("%s/orders/%s?email=%s", p.cfg.URL, orderHash, email))
	mq := messenger.MessageQuery{}
	mq.RecipientID(to)
	mq.Template(template.ButtonTemplate{
		Buttons:      []template.Button{orderDetailsBtn},
		Text:         msg,
		TemplateType: template.TemplateTypeButton,
	})

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return err
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)
	return nil
}

func (p *FacebookStateProcessor) processPostback(event messenger.Event, opts messenger.MessageOpts, postback messenger.Postback) error {
	log.Log().Infoln("Message processPostback....")

	if strings.Contains(strings.ToLower(postback.Payload), "say, hello") {
		p.sendWelcomeMessage(opts.Sender.ID)
		return nil
	}

	foundAction, action, args := p.findActionByPattern(postback.Payload)
	if !foundAction {
		p.sendDefaultMessage(opts.Sender.ID)
		return nil
	}
	action.Fn(opts.Sender.ID, action.MenuTemplate, p.parseArgs(args))
	return nil
}

func (p *FacebookStateProcessor) processMessage(event messenger.Event, opts messenger.MessageOpts, message messenger.ReceivedMessage) error {
	if message.QuickReply != nil {
		return p.processQuickReply(event, opts, message.QuickReply)
	}

	log.Log().Infoln("Message processMessage....")

	foundAction, action, args := p.findActionByPattern(message.Text)
	if !foundAction {
		customerState, _ := p.stateService.GetState(opts.Sender.ID)
		switch customerState {
		case models.CustomerStateSearchProducts:
			args := []string{"1", message.Text}
			p.searchProducts(opts.Sender.ID, template.MenuTemplateSearchProductsWithQuery, args)
		default:
			if strings.Contains(strings.ToLower(message.Text), "say, hello") {
				p.sendWelcomeMessage(opts.Sender.ID)
			} else {
				p.sendDefaultMessage(opts.Sender.ID)
			}
		}
		return nil
	}
	action.Fn(opts.Sender.ID, action.MenuTemplate, p.parseArgs(args))
	return nil
}

func (p *FacebookStateProcessor) processQuickReply(event messenger.Event, opts messenger.MessageOpts, message *messenger.QuickReplyPayload) error {
	log.Log().Infoln("Message processQuickReply....")

	foundAction, action, args := p.findActionByPattern(message.Payload)
	if !foundAction {
		p.sendDefaultMessage(opts.Sender.ID)
		return nil
	}
	action.Fn(opts.Sender.ID, action.MenuTemplate, p.parseArgs(args))
	return nil
}

func (p *FacebookStateProcessor) loadMenus() []template.Button {
	var menus []template.Button
	for _, a := range p.allActions() {
		if a.IsMenu {
			menus = append(menus, template.NewPostbackButton(a.Name, a.MenuTemplate))
		}
	}
	return menus
}

func (p *FacebookStateProcessor) addToCart(senderID, tmpl string, args []string) {
	log.Log().Infoln("Requesting addToCart | userID: ", senderID, " | args: ", args)

	cartID, err := p.stateService.GetData(senderID, "cart")
	if err != nil {
		log.Log().Errorln(err)
		return
	}

	quantity := 1
	if len(args) > 0 {
		q, err := strconv.ParseInt(args[0], 10, 32)
		if err == nil {
			quantity = int(q)
			if quantity < 0 {
				quantity = 0
			}
		}
	}
	productID := args[1]

	var cart *models.Cart

	if cartID == nil {
		cart, err = p.shopemaaService.CreateCart(productID, quantity)
		if err != nil {
			if !strings.Contains(err.Error(), "out of stock") {
				log.Log().Errorln(err)
				return
			}

			resp, err := p.messenger.SendSimpleMessage(senderID, "Out of stock")
			if err != nil {
				log.Log().Errorln(err)
				return
			}
			log.Log().Infoln("MessageID: ", resp.MessageID)
			return
		}
	} else {
		cart, err = p.shopemaaService.UpdateCart(cartID.(string), productID, quantity)
		if err != nil {
			if !strings.Contains(err.Error(), "out of stock") {
				log.Log().Errorln(err)
				return
			}

			resp, err := p.messenger.SendSimpleMessage(senderID, "Out of stock")
			if err != nil {
				log.Log().Errorln(err)
				return
			}
			log.Log().Infoln("MessageID: ", resp.MessageID)
			return
		}
	}

	p.stateService.SetData(senderID, "cart", cart.ID)

	log.Log().Infoln("Cart updated. Items: ", len(cart.CartItems))

	if len(cart.CartItems) == 0 {
		if _, err := p.messenger.SendSimpleMessage(senderID, "Cart is empty."); err != nil {
			log.Log().Errorln(err)
		}
		return
	}

	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)

	for _, p := range cart.CartItems {
		imageUrl := ""
		if len(p.Product.FullImages) > 0 {
			imageUrl = p.Product.FullImages[0]
		}
		mq.Template(template.GenericTemplate{
			Title:    fmt.Sprintf("%s", p.Product.Name),
			Subtitle: fmt.Sprintf("BDT %0.2f | Qty: %d", float64(p.PurchasePrice)/float64(100), p.Quantity),
			ImageURL: imageUrl,
			Buttons: []template.Button{
				{
					Title:   "+",
					Payload: fmt.Sprintf(template.MenuTemplateProductAddToCart, p.Product.ID, p.Quantity+1),
					Type:    template.ButtonTypePostback,
				},
				{
					Title:   "-",
					Payload: fmt.Sprintf(template.MenuTemplateProductAddToCart, p.Product.ID, p.Quantity-1),
					Type:    template.ButtonTypePostback,
				},
			},
		})
	}

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)

	p.sendCheckoutOption(senderID, cart.ID)
}

func (p *FacebookStateProcessor) sendCheckoutOption(userID, cartID string) {
	continueShoppingBtn := template.NewPostbackButton("Continue", template.MenuTemplateSearchProductsRequestDefault)
	checkoutBtn := template.NewWebURLButton("Checkout", fmt.Sprintf("%s/checkout/%s", p.cfg.URL, cartID))
	mq := messenger.MessageQuery{}
	mq.RecipientID(userID)
	mq.Template(template.ButtonTemplate{
		Buttons:      []template.Button{continueShoppingBtn, checkoutBtn},
		Text:         "What's next?",
		TemplateType: template.TemplateTypeButton,
	})

	p.stateService.SetIdentityByCartID("to", cartID, userID)

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) showCart(senderID, tmpl string, args []string) {
	log.Log().Infoln("Requesting showCart | userID: ", senderID)

	cartID, err := p.stateService.GetData(senderID, "cart")
	if err != nil {
		log.Log().Errorln(err)
		return
	}

	var cart *models.Cart

	if cartID == nil {
		if _, err := p.messenger.SendSimpleMessage(senderID, "Cart is empty."); err != nil {
			log.Log().Errorln(err)
		}
		return
	} else {
		cart, err = p.shopemaaService.GetCart(cartID.(string))
		if err != nil {
			log.Log().Errorln(err)
			return
		}
	}

	if len(cart.CartItems) == 0 {
		if _, err := p.messenger.SendSimpleMessage(senderID, "Cart is empty."); err != nil {
			log.Log().Errorln(err)
		}
		return
	}

	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)

	for _, p := range cart.CartItems {
		imageUrl := ""
		if len(p.Product.FullImages) > 0 {
			imageUrl = p.Product.FullImages[0]
		}
		mq.Template(template.GenericTemplate{
			Title:    fmt.Sprintf("%s", p.Product.Name),
			Subtitle: fmt.Sprintf("BDT %0.2f | Qty: %d", float64(p.PurchasePrice)/float64(100), p.Quantity),
			ImageURL: imageUrl,
			Buttons: []template.Button{
				{
					Title:   "+",
					Payload: fmt.Sprintf(template.MenuTemplateProductAddToCart, p.Product.ID, p.Quantity+1),
					Type:    template.ButtonTypePostback,
				},
				{
					Title:   "-",
					Payload: fmt.Sprintf(template.MenuTemplateProductAddToCart, p.Product.ID, p.Quantity-1),
					Type:    template.ButtonTypePostback,
				},
			},
		})
	}

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)

	p.sendCheckoutOption(senderID, cart.ID)
}

func (p *FacebookStateProcessor) talkToUs(senderID, tmpl string, args []string) {
	log.Log().Infoln("Requesting talkToUs | userID: ", senderID)

	callUsBtn := template.NewPhoneNumberButton("Call Us", p.shopemaaService.GetShop().SupportPhone)
	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)
	mq.Template(template.ButtonTemplate{
		Buttons:      []template.Button{callUsBtn},
		Text:         "Have a query?",
		TemplateType: template.TemplateTypeButton,
	})

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) searchProducts(senderID, tmpl string, args []string) {
	log.Log().Infoln("searchProducts triggered...userID: ", senderID, " | args: ", args)

	limit := 5
	currentPage := 1
	if len(args) > 0 {
		page, err := strconv.ParseInt(args[0], 10, 32)
		if err == nil {
			currentPage = int(page)
			if currentPage < 1 {
				currentPage = 1
			}
		}
	}

	query := ""
	if len(args) >= 2 {
		query = args[1]
	}

	products, err := p.shopemaaService.SearchProducts(query, currentPage, limit)
	if err != nil {
		log.Log().Errorln(err)
		return
	}

	if len(products) == 0 {
		if currentPage == 1 {
			if _, err := p.messenger.SendSimpleMessage(senderID, "No product found."); err != nil {
				log.Log().Errorln(err)
			}
		} else {
			if _, err := p.messenger.SendSimpleMessage(senderID, "No more products."); err != nil {
				log.Log().Errorln(err)
			}
		}
		return
	}

	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)
	for _, prod := range products {
		imageUrl := ""
		if len(prod.FullImages) > 0 {
			imageUrl = prod.FullImages[0]
		}
		mq.Template(template.GenericTemplate{
			Title:    fmt.Sprintf("%s", prod.Name),
			Subtitle: fmt.Sprintf("%s %0.2f", p.shopemaaService.GetCurrency(), float64(prod.Price)/float64(100)),
			ImageURL: imageUrl,
			Buttons: []template.Button{
				{
					Title:   "Add to Cart",
					Payload: fmt.Sprintf(template.MenuTemplateProductAddToCart, prod.ID, 1),
					Type:    template.ButtonTypePostback,
				},
			},
		})
	}

	mq.QuickReply(messenger.QuickReply{
		Title:       "Load More",
		Payload:     fmt.Sprintf(template.MenuTemplateSearchProductsWithQuery, query, currentPage+1),
		ContentType: messenger.ContentTypeText,
	})

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) searchProductsRequest(senderID, template string, args []string) {
	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)
	mq.Text("What are you looking for?")
	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	p.stateService.SetState(senderID, models.CustomerStateSearchProducts)
	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) listCategories(senderID, tmpl string, args []string) {
	log.Log().Infoln("Args: ", args)

	limit := 5
	currentPage := 1
	if len(args) > 0 {
		page, err := strconv.ParseInt(args[0], 10, 32)
		if err == nil {
			currentPage = int(page)
			if currentPage < 1 {
				currentPage = 1
			}
		}
	}

	categories, err := p.shopemaaService.ListCategories(currentPage, limit)
	if err != nil {
		log.Log().Errorln(err)
		return
	}

	if len(categories) == 0 {
		if _, err := p.messenger.SendSimpleMessage(senderID, "No more categories."); err != nil {
			log.Log().Errorln(err)
		}
		return
	}

	mq := messenger.MessageQuery{}
	mq.Text("Categories")
	mq.RecipientID(senderID)
	for _, c := range categories {
		mq.QuickReply(messenger.QuickReply{
			Title:       fmt.Sprintf("%s (%d)", c.Name, c.ProductCount),
			Payload:     fmt.Sprintf(template.MenuTemplateSearchProductsByCategory, c.ID, 1),
			ContentType: messenger.ContentTypeText,
		})
	}
	mq.QuickReply(messenger.QuickReply{
		Title:       "Prev",
		Payload:     fmt.Sprintf(tmpl, currentPage-1),
		ContentType: messenger.ContentTypeText,
	})
	mq.QuickReply(messenger.QuickReply{
		Title:       "Next",
		Payload:     fmt.Sprintf(tmpl, currentPage+1),
		ContentType: messenger.ContentTypeText,
	})

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}

	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) searchProductsByCategories(senderID, tmpl string, args []string) {
	log.Log().Infoln("searchProductsByCategories triggered...userID: ", senderID, " | args: ", args)

	limit := 5
	currentPage := 1
	if len(args) > 0 {
		page, err := strconv.ParseInt(args[0], 10, 32)
		if err == nil {
			currentPage = int(page)
			if currentPage < 1 {
				currentPage = 1
			}
		}
	}

	categoryID := ""
	if len(args) >= 2 {
		categoryID = args[1]
	}

	products, err := p.shopemaaService.ListProductsByCategory(categoryID, currentPage, limit)
	if err != nil {
		log.Log().Errorln(err)
		return
	}

	if len(products) == 0 {
		if currentPage == 1 {
			if _, err := p.messenger.SendSimpleMessage(senderID, "No product found."); err != nil {
				log.Log().Errorln(err)
			}
		} else {
			if _, err := p.messenger.SendSimpleMessage(senderID, "No more products."); err != nil {
				log.Log().Errorln(err)
			}
		}
		return
	}

	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)
	for _, p := range products {
		imageUrl := ""
		if len(p.FullImages) > 0 {
			imageUrl = p.FullImages[0]
		}
		mq.Template(template.GenericTemplate{
			Title:    fmt.Sprintf("%s", p.Name),
			Subtitle: fmt.Sprintf("BDT %0.2f", float64(p.Price)/float64(100)),
			ImageURL: imageUrl,
			Buttons: []template.Button{
				{
					Title:   "Add to Cart",
					Payload: fmt.Sprintf(template.MenuTemplateProductAddToCart, p.ID, 1),
					Type:    template.ButtonTypePostback,
				},
			},
		})
	}

	mq.QuickReply(messenger.QuickReply{
		Title:       "Load More",
		Payload:     fmt.Sprintf(template.MenuTemplateSearchProductsByCategory, categoryID, currentPage+1),
		ContentType: messenger.ContentTypeText,
	})

	resp, err := p.messenger.SendMessage(mq)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) sendDefaultMessage(senderID string) {
	msg := "Jally here! I don't understand human language, send me commands!"
	resp, err := p.messenger.SendSimpleMessage(senderID, msg)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) sendWelcomeMessage(senderID string) {
	msg := fmt.Sprintf("Welcome to %s!!", p.shopemaaService.GetName())
	resp, err := p.messenger.SendSimpleMessage(senderID, msg)
	if err != nil {
		log.Log().Errorln(err)
		return
	}

	msg = "What are you looking for?"
	resp, err = p.messenger.SendSimpleMessage(senderID, msg)
	if err != nil {
		log.Log().Errorln(err)
		return
	}
	p.stateService.SetState(senderID, models.CustomerStateSearchProducts)
	log.Log().Infoln("MessageID: ", resp.MessageID)
}

func (p *FacebookStateProcessor) allActions() []Action {
	return []Action{
		{Name: "Categories", Fn: p.listCategories, MenuPattern: template.MenuPatternCategories, MenuTemplate: template.MenuTemplateCategoriesDefault, IsMenu: true},
		{Name: "Search Products", Fn: p.searchProductsRequest, MenuPattern: template.MenuPatternSearchProductsRequest, MenuTemplate: template.MenuTemplateSearchProductsRequestDefault, IsMenu: true},
		{Name: "My Cart", Fn: p.showCart, MenuPattern: template.MenuPatternShowMyCart, MenuTemplate: template.MenuPatternShowMyCart, IsMenu: true},
		{Name: "Talk to Us", Fn: p.talkToUs, MenuPattern: template.MenuPatternTalkToUs, MenuTemplate: template.MenuPatternTalkToUs, IsMenu: true},
		{Name: "Search Products", Fn: p.searchProducts, MenuPattern: template.MenuPatternSearchProducts, MenuTemplate: template.MenuTemplateSearchProducts},
		{Name: "Search Products With Query", Fn: p.searchProducts, MenuPattern: template.MenuPatternSearchProductsWithQuery, MenuTemplate: template.MenuTemplateSearchProductsWithQuery},
		{Name: "List Products by Category", Fn: p.searchProductsByCategories, MenuPattern: template.MenuPatternSearchProductsByCategory, MenuTemplate: template.MenuTemplateSearchProductsByCategory},
		{Name: "Add to Cart", Fn: p.addToCart, MenuPattern: template.MenuPatternProductAddToCart, MenuTemplate: template.MenuTemplateProductAddToCart},
	}
}

func (p *FacebookStateProcessor) parseArgs(args []string) []string {
	log.Log().Infoln("Args: ", args)

	var parsed []string

	for _, a := range args {
		a = strings.ReplaceAll(a, "search_", "")
		a = strings.ReplaceAll(a, "products_", "")
		a = strings.ReplaceAll(a, "with_", "")
		a = strings.ReplaceAll(a, "by_", "")
		a = strings.ReplaceAll(a, "categories_", "")
		a = strings.ReplaceAll(a, "category_", "")
		a = strings.ReplaceAll(a, "product_add_to_cart_", "")

		// Parsing Int
		re, _ := regexp.Compile("_[0-9]+")
		values := re.FindAllString(a, -1)
		for _, v := range values {
			parsed = append(parsed, strings.Trim(v, "_"))
		}

		// Parsing UUID
		re, _ = regexp.Compile("([a-zA-Z0-9-]{36})+")
		values = re.FindAllString(a, -1)
		for _, v := range values {
			parsed = append(parsed, strings.Trim(v, "_"))
		}

		// Parsing Query
		re, _ = regexp.Compile("[a-zA-Z0-9 ]+")
		values = re.FindAllString(a, -1)
		for _, v := range values {
			parsed = append(parsed, strings.Trim(v, "_"))
		}
	}

	return parsed
}

func (p *FacebookStateProcessor) findActionByPattern(pattern string) (bool, *Action, []string) {
	log.Log().Infoln("Pattern: " + pattern)

	for _, a := range p.allActions() {
		r, _ := regexp.Compile(a.MenuPattern)
		matches := r.FindAllString(pattern, -1)
		if len(matches) > 0 {
			return true, &a, matches
		}
	}
	return false, nil, nil
}

func NewFacebookStateProcessor(cfg *config.Application, stateService services.IStateService,
	shopemaaService services.IShopemaaService, messenger *messenger.Messenger) (IStateProcessor, error) {
	processor := &FacebookStateProcessor{
		cfg:             cfg,
		stateService:    stateService,
		shopemaaService: shopemaaService,
		messenger:       messenger,
	}
	if err := processor.Init(); err != nil {
		return nil, err
	}
	return processor, nil
}
