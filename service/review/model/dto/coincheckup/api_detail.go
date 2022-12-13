package dto_coincheckup

//https://coincheckup.com/api/coincheckup/get_coin_assets_byslug/uniswap

//Crawl:
//	Anh san pham
//	exchange info cao them o url: https://coincheckup.com/coins/{Market.CoinId}/markets map voi lowercase exchange id
//
type APIDetail struct {
	Category    []string `json:"categories"`
	Exchanges   []string `json:"exchanges"` //slug(id) in many other website
	Description string   `json:"description"`
	MarketData  Market   `json:"market"`
	SocialsData []Social `json:"socials"`
	Website     string   `json:"website"`
	WhitePaper  string   `json:"whitepaper"`
	////////////////////////////////

	Research Research `json:"research"` //busd, alias don't have this data
}

//`n/a` or `` for null
type Research struct {
	LatestUpdated string `json:"research_update_date"`

	//Website
	WebsiteUrl       string `json:"website_url"`
	Website2Url      string `json:"website2_url"`
	BlogUrl          string `json:"blog_url"`
	FaqUrl           string `json:"faq_url"`
	MessageBoard1Url string `json:"message_board_1_url"`
	MessageBoard2Url string `json:"message_board_2_url"`
	MessageBoard3Url string `json:"message_board_3_url"`
	BitcointalkUrl   string `json:"bitcointalk_url"`
	ChainExp1Url     string `json:"chain_exp_1_url"`
	ChainExp2Url     string `json:"chain_exp_2_url"`
	ChainExp3Url     string `json:"chain_exp_3_url"`
	WhitepaperUrl    string `json:"whitepaper_url"`
	FacebookUrl      string `json:"facebook_url"`
	TwitterUrl       string `json:"twitter_url"`
	YoutubeUrl       string `json:"youtube_url"`
	SlackUrl         string `json:"slack_url"`
	TelegramUrl      string `json:"telegram_url"`
	IcqHandle        string `json:"icq_handle"`
	GooglePlusUrl    string `json:"google_plus_url"`
	OtherSocialUrl   string `json:"other_social_url"`
	GithubUrl        string `json:"github_url"`
	TeamPageUrl      string `json:"team_page_url"`
	AdvantageUrl     string `json:"company_explains_advantage_over_competition_source"`
	MktPlanUrl       string `json:"company_presents_sales_mkt_plan_source"`
	RedditUrl        string `json:"reddit_url"`

	//Contact
	ContactEmail string `json:"contact_email"`
	PhoneNum     string `json:"phone_num"`

	//Description
	CoinGoals                 string   `json:"coin_goals"`
	AdditionalDifferentiation string   `json:"additional_differentiation"`
	Video                     string   `json:"video_what_is_it"`
	DescriptionShort          string   `json:"description_short"`
	Abstract                  string   `json:"abstract"`
	ProductRoadmapUrl         string   `json:"product_roadmap_url"`
	Tags                      []string `json:"categories"`

	//Technique Info
	InternalSignature          string  `json:"internal_signature"`
	DomainStartDate            string  `json:"whois_reg_date"`
	Minable                    string  `json:"minable"`
	Algorithm                  string  `json:"algorithm"`
	TokenSupply                float64 `json:"token_supply"` //zero is no limit
	ConsensusMethod            string  `json:"consensus_method"`
	MinuteBlockProcessingSpeed string  `json:"block_processing_speed"` //Minnute per block
	TransactionsPerSecond      float64 `json:"transactions_per_second"`
	Governance                 string  `json:"governance"`
	ProductStatus              string  `json:"product_status"`
	ProductReleaseDate         string  `json:"product_release_date"`

	//Manager Info
	TeamDevSize         string `json:"team_dev_size"`
	TeamMktSalesSize    string `json:"team_mkt_sales_size"`
	TeamSizeTotal       string `json:"team_size_total"`
	TeamAgeAvg          string `json:"team_age_avg"`
	CeoName             string `json:"ceo_name"`
	CeoLinkedinUrl      string `json:"ceo_linkedin_url"`
	CeoPriorEngagements string `json:"ceo_prior_engagements"`
	CtoName             string `json:"cto_name"`
	CtoLinkedinUrl      string `json:"cto_linkedin_url"`
	CtoGithubUrl        string `json:"cto_github_url"`
	CtoPriorEngagements string `json:"cto_prior_engagements"`

	//ICO Info
	IcoProceedsFoundersTeam  string `json:"ico_proceeds_founders_team"`
	IcoProceedsMktSales      string `json:"ico_proceeds_mkt_sales"`
	IcoProceedsDevelopment   string `json:"ico_proceeds_development"`
	IcoProceedsAdmLegal      string `json:"ico_proceeds_adm_legal"`
	IcoProceedsOther         string `json:"ico_proceeds_other"`
	IcoSalePeriodStart       string `json:"ico_sale_period_start"`
	IcoSalePeriodEnd         string `json:"ico_sale_period_end"`
	IcoInitialPrice          string `json:"ico_initial_price"`
	IcoAcceptedCurrencies    string `json:"ico_accepted_currencies"`
	IcoInvestmentRound       string `json:"ico_investment_round"`
	IcoTokenDistributionDate string `json:"ico_token_distribution_date"`
	IcoTokenTradingStartDate string `json:"ico_token_trading_start_date"`
	IcoMinInvestmentGoal     string `json:"ico_min_investment_goal"`
	IcoMaxInvestmentCap      string `json:"ico_max_investment_cap"`
	IcoBonus                 string `json:"ico_bonus"`
	IcoJurisdiction          string `json:"ico_jurisdiction"`
	IcoLegalAdvisers         string `json:"ico_legal_advisers"`
	IcoLegalForm             string `json:"ico_legal_form"`
	IcoSecurity              string `json:"ico_security"`
}

type Market struct {
	IcoPrice    any         `json:"ico_price"` // float64
	InitialData InitialData `json:"initial_data"`
	CoinSymbol  string      `json:"symbol"`
	CoinName    string      `json:"name"`
	CoinId      string      `json:"ccx_slug"`
}

type InitialData struct {
	PriceUsd string `json:"price_usd"`
	Date     string `json:"date"`
}

type Social struct {
	Name                string `json:"name"`
	Id                  string `json:"id"`
	CoincodexCoinSymbol string `json:"coincodex_coin_symbol"`
	CoincodexSocialsId  string `json:"coincodex_socials_id"`
	Value               string `json:"value"`
	Label               string `json:"label"`
	OrderBy             string `json:"order_by"`
}
