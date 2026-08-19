package landing

// HeroSection data
type HeroSection struct {
	Badge       string `json:"badge"`
	Headline    string `json:"headline"`
	Subheadline string `json:"subheadline"`
	PrimaryCTA  string `json:"primary_cta"`
	SecondaryCTA string `json:"secondary_cta"`
	RatingNote  string `json:"rating_note"`
}

// FeatureItem data
type FeatureItem struct {
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// HowItWorksStep data
type HowItWorksStep struct {
	Step        int    `json:"step"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// FAQItem data
type FAQItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// ReviewItem data
type ReviewItem struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Rating  int    `json:"rating"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

// TemplatePreviewItem data
type TemplatePreviewItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Category    string `json:"category"`
	ThumbnailURL string `json:"thumbnail_url"`
	IsFeatured  bool   `json:"is_featured"`
}

// LandingData merepresentasikan struktur data lengkap landing page yang disajikan ke guest.
// Sesuai SSOT §4.1 (Landing Page Guest), §43 (Landing Page Management), §66 (Landing & Package API).
type LandingData struct {
	Hero         HeroSection           `json:"hero"`
	Features     []FeatureItem         `json:"features"`
	HowItWorks   []HowItWorksStep      `json:"how_it_works"`
	Templates    []TemplatePreviewItem `json:"templates"`
	Reviews      []ReviewItem          `json:"reviews"`
	FAQs         []FAQItem             `json:"faqs"`
	ContactInfo  map[string]string     `json:"contact_info"`
}
