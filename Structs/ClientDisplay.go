package Structs

// Deploying to app store is a pain and taking down the client web service can impact user experience.
// This struct is for the client display info to make it easier for the frontend guys.
type ClientDisplay struct {
	PageTitle string `json:"pageTitle"`

	//Android, Web, and IOS have different elements. We can get platform in request and get display settings.
	PageElements []PageElement `json:"pageElements"`
}

// For each element on the page, if not given client auto assumes
// Must schedule meeting with client side guys before approval.
type PageElement struct {
	Id                 string  `json:"id"` // Element ID or client-specific identifier
	BackgroundColor    string  `json:"backgroundColor,omitempty"`
	BackgroundImage    string  `json:"backgroundImage,omitempty"`
	TextSize           string  `json:"textSize,omitempty"`           // Font size
	TextWeight         string  `json:"textWeight,omitempty"`         // Font weight (e.g., bold, normal)
	TextFont           string  `json:"textFont,omitempty"`           // Font family
	TextHighlightColor string  `json:"textHighlightColor,omitempty"` // Text highlight color
	TextColor          string  `json:"textColor,omitempty"`          // Text color
	FontStyle          string  `json:"fontStyle,omitempty"`          // Font style (e.g., italic)
	TextDecoration     string  `json:"textDecoration,omitempty"`     // Text decoration (e.g., underline)
	LetterSpacing      string  `json:"letterSpacing,omitempty"`      // Letter spacing
	LineHeight         string  `json:"lineHeight,omitempty"`         // Line height
	TextTransform      string  `json:"textTransform,omitempty"`      // Text transformation (e.g., uppercase)
	Border             string  `json:"border,omitempty"`             // Border properties
	BorderRadius       string  `json:"borderRadius,omitempty"`       // Border radius
	BoxShadow          string  `json:"boxShadow,omitempty"`          // Box shadow
	Opacity            float64 `json:"opacity,omitempty"`            // Opacity level
	Width              string  `json:"width,omitempty"`              // Element width
	Height             string  `json:"height,omitempty"`             // Element height
	Margin             string  `json:"margin,omitempty"`             // Margin
	Padding            string  `json:"padding,omitempty"`            // Padding
	Display            string  `json:"display,omitempty"`            // Display type (e.g., block, inline)
	Position           string  `json:"position,omitempty"`           // Positioning (e.g., relative, absolute)
	Top                string  `json:"top,omitempty"`                // Top position
	Bottom             string  `json:"bottom,omitempty"`             // Bottom position
	Left               string  `json:"left,omitempty"`               // Left position
	Right              string  `json:"right,omitempty"`              // Right position
	ZIndex             int     `json:"zIndex,omitempty"`             // Z-index
	Overflow           string  `json:"overflow,omitempty"`           // Overflow behavior
	TextAlign          string  `json:"textAlign,omitempty"`          // Text alignment (e.g., left, center, right)
	VerticalAlign      string  `json:"verticalAlign,omitempty"`      // Vertical alignment
	BackgroundSize     string  `json:"backgroundSize,omitempty"`     // Background size
	BackgroundPosition string  `json:"backgroundPosition,omitempty"` // Background position
	BackgroundRepeat   string  `json:"backgroundRepeat,omitempty"`   // Background repeat behavior
	Cursor             string  `json:"cursor,omitempty"`             // Cursor type
	Transition         string  `json:"transition,omitempty"`         // Transition effects
	Animation          string  `json:"animation,omitempty"`          // Animation effects

	Float           string `json:"float,omitempty"`           // Floating behavior (e.g., left, right)
	Clear           string `json:"clear,omitempty"`           // Clearing behavior (e.g., left, right, both)
	Visibility      string `json:"visibility,omitempty"`      // Element visibility (e.g., visible, hidden)
	PointerEvents   string `json:"pointerEvents,omitempty"`   // Pointer events behavior (e.g., auto, none)
	UserSelect      string `json:"userSelect,omitempty"`      // User selection behavior (e.g., none, text)
	BoxSizing       string `json:"boxSizing,omitempty"`       // Box sizing model (e.g., border-box, content-box)
	Outline         string `json:"outline,omitempty"`         // Outline properties
	Resize          string `json:"resize,omitempty"`          // Resizability (e.g., both, horizontal, vertical)
	Filter          string `json:"filter,omitempty"`          // Image effects (e.g., blur, grayscale)
	ImageRendering  string `json:"imageRendering,omitempty"`  // Image rendering quality (e.g., auto, crisp-edges)
	ShapeOutside    string `json:"shapeOutside,omitempty"`    // CSS Shapes outside property
	ClipPath        string `json:"clipPath,omitempty"`        // Clip path property
	BackdropFilter  string `json:"backdropFilter,omitempty"`  // Backdrop filter effects
	MixBlendMode    string `json:"mixBlendMode,omitempty"`    // Mix blend mode for element blending
	Transform       string `json:"transform,omitempty"`       // Transformations (e.g., rotate, scale)
	Perspective     string `json:"perspective,omitempty"`     // 3D perspective effect
	TransformOrigin string `json:"transformOrigin,omitempty"` // Transform origin point
	WillChange      string `json:"willChange,omitempty"`      // Property hinting for performance optimization
	ScrollBehavior  string `json:"scrollBehavior,omitempty"`  // Scroll behavior (e.g., smooth)

	// Custom properties for specific use cases
	CustomCSS string `json:"customCSS,omitempty"` // Additional custom CSS
}
