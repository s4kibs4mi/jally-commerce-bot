package template

const (
	MenuPatternCategories         string = "categories_[1-9]+"
	MenuTemplateCategories        string = "categories_%d"
	MenuTemplateCategoriesDefault string = "categories_1"

	MenuPatternSearchProductsRequest         string = "search_products_request"
	MenuTemplateSearchProductsRequest        string = "search_products_request"
	MenuTemplateSearchProductsRequestDefault string = "search_products_request"

	MenuPatternSearchProducts           string = "search_products_[1-9]+"
	MenuTemplateSearchProducts          string = "search_products_%d"
	MenuPatternSearchProductsWithQuery  string = "search_products_with_[a-zA-Z0-9 ]+_[1-9]+"
	MenuTemplateSearchProductsWithQuery string = "search_products_with_%s_%d"

	MenuPatternProductAddToCart  string = "product_add_to_cart_[a-zA-Z0-9-]+_[0-9]+"
	MenuTemplateProductAddToCart string = "product_add_to_cart_%s_%d"

	MenuPatternSearchProductsByCategory  string = "search_products_by_category_[a-zA-Z0-9-]+_[1-9]+"
	MenuTemplateSearchProductsByCategory string = "search_products_by_category_%s_%d"

	MenuPatternShowMyCart string = "show_my_cart"
	MenuPatternTalkToUs   string = "talk_to_us"
)
