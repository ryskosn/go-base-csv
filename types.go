package main

// Orders is slice of Order
// https://github.com/baseinc/api-docs/blob/master/base_api_v1_orders.md
type Orders struct {
	Orders []struct {
		UniqueKey      string      `json:"unique_key"`
		Ordered        int         `json:"ordered"`
		Cancelled      interface{} `json:"cancelled"`
		Dispatched     int         `json:"dispatched"`
		Payment        string      `json:"payment"`
		FirstName      string      `json:"first_name"`
		LastName       string      `json:"last_name"`
		Total          int         `json:"total"`
		Terminated     bool        `json:"terminated"`
		DispatchStatus string      `json:"dispatch_status"`
		Modified       int         `json:"modified"`
	} `json:"orders"`
}

// Order is
// https://github.com/baseinc/api-docs/blob/master/base_api_v1_orders_detail.md
type Order struct {
	Order struct {
		UniqueKey         string      `json:"unique_key"`
		Ordered           int         `json:"ordered"`
		Cancelled         interface{} `json:"cancelled"`
		Dispatched        interface{} `json:"dispatched"`
		Payment           string      `json:"payment"`
		ShippingMethod    string      `json:"shipping_method"`
		ShippingFee       int         `json:"shipping_fee"`
		CodFee            int         `json:"cod_fee"`
		Total             int         `json:"total"`
		FirstName         string      `json:"first_name"`
		LastName          string      `json:"last_name"`
		Country           string      `json:"country"`
		ZipCode           string      `json:"zip_code"`
		Prefecture        string      `json:"prefecture"`
		Address           string      `json:"address"`
		Address2          string      `json:"address2"`
		MailAddress       string      `json:"mail_address"`
		Tel               string      `json:"tel"`
		Remark            string      `json:"remark"`
		AddComment        string      `json:"add_comment"`
		DeliveryCompanyID int         `json:"delivery_company_id"`
		TrackingNumber    string      `json:"tracking_number"`
		Terminated        bool        `json:"terminated"`
		DispatchStatus    string      `json:"dispatch_status"`
		Modified          int         `json:"modified"`
		OrderReceiver     struct {
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			ZipCode    string `json:"zip_code"`
			Prefecture string `json:"prefecture"`
			Address    string `json:"address"`
			Address2   string `json:"address2"`
			Tel        string `json:"tel"`
		} `json:"order_receiver"`
		OrderDiscount struct {
			Discount int    `json:"discount"`
			Note     string `json:"note"`
		} `json:"order_discount"`
		CCPaymentTransaction struct {
			CollectedFee int `json:"collected_fee"`
		} `json:"c_c_payment_transaction"`
		CvsPaymentTransaction struct {
			CollectedFee interface{} `json:"collected_fee"`
			Status       interface{} `json:"status"`
		} `json:"cvs_payment_transaction"`
		BtPaymentTransaction struct {
			CollectedFee interface{} `json:"collected_fee"`
			Status       interface{} `json:"status"`
		} `json:"bt_payment_transaction"`
		AtobaraiPaymentTransaction struct {
			CollectedFee interface{} `json:"collected_fee"`
			Status       interface{} `json:"status"`
		} `json:"atobarai_payment_transaction"`
		OrderItems []struct {
			OrderItemID    int         `json:"order_item_id"`
			ItemID         int         `json:"item_id"`
			VariationID    int         `json:"variation_id"`
			Title          string      `json:"title"`
			Variation      string      `json:"variation"`
			Price          int         `json:"price"`
			Amount         int         `json:"amount"`
			Total          int         `json:"total"`
			Status         string      `json:"status"`
			ShippingMethod interface{} `json:"shipping_method"`
			ShippingFee    int         `json:"shipping_fee"`
			Modified       int         `json:"modified"`
		} `json:"order_items"`
	} `json:"order"`
}

// OrderToCsv is struct of output csv format
type OrderToCsv struct {
	UniqueKey   string `json:"unique_key"`
	Ordered     string `json:"ordered"`
	FullName    string
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	Total       int    `json:"total"`
	MailAddress string `json:"mail_address"`
	Tel         string `json:"tel"`
	Remark      string `json:"remark"`
	Discount    int    `json:"discount"`
	Coupon      string
	Note        string `json:"note"`
	ItemID      int    `json:"item_id"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Country     string `json:"country"`
	ZipCode     string `json:"zip_code"`
	Prefecture  string `json:"prefecture"`
	// Address     string `json:"address"`
	// Address2    string `json:"address2"`

	// Cancelled      interface{} `json:"cancelled"`
	// Dispatched     interface{} `json:"dispatched"`
	// Payment        string      `json:"payment"`
	// DispatchStatus string      `json:"dispatch_status"`
}

// Items is
// https://github.com/baseinc/api-docs/blob/master/base_api_v1_items.md
type Items struct {
	Items []struct {
		ItemID      int         `json:"item_id"`
		Title       string      `json:"title"`
		Detail      string      `json:"detail"`
		Price       int         `json:"price"`
		ProperPrice interface{} `json:"proper_price"`
		Stock       int         `json:"stock"`
		Visible     int         `json:"visible"`
		ListOrder   int         `json:"list_order"`
		Identifier  string      `json:"identifier"`
		Img1Origin  string      `json:"img1_origin"`
		Img2Origin  string      `json:"img2_origin"`
		Img3Origin  interface{} `json:"img3_origin"`
		Img4Origin  interface{} `json:"img4_origin"`
		Img5Origin  interface{} `json:"img5_origin"`
		Modified    int         `json:"modified"`
		Variations  []struct {
			VariationID         int    `json:"variation_id"`
			Variation           string `json:"variation"`
			VariationStock      int    `json:"variation_stock"`
			VariationIdentifier string `json:"variation_identifier"`
		} `json:"variations"`
	} `json:"items"`
}
