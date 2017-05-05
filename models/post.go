package models

// CorePost struct to include info of core Post
type CorePost struct {
	UserObject
	PostID    int    `json:"id"` //Id
	Message   string `json:"message"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	Status    int    `json:"status,omitempty"`
	Privacy   int    `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
	Likes     int    `json:"likes,omitempty"`
	Comments  int    `json:"comments,omitempty"`
	Shares    int    `json:"shares,omitempty"`
	IsLiked   bool   `json:"is_liked"`
	Type      int    `json:"type"`         // 1: wall post; 2: group post; 3: upload photo; 4: share;
	To        Object `json:"to,omitempty"` // include objectid, name
}

// WallPost struct
type WallPost struct {
	CorePost
	Link string `json:"link,omitempty"`
}

// PhotoPost struct
type PhotoPost struct {
	CorePost
	Photo string `json:"photo,omitempty"`
}

// LinkPost struct
type LinkPost struct {
	CorePost
	Link string `json:"link,omitempty"`
}

// SharePost struct
type SharePost struct {
	CorePost
	Content interface{} `json:"content,omitempty"`
}

// Post struct
type Post struct {
	CorePost
	Link    string      `json:"link,omitempty"`
	Content interface{} `json:"content,omitempty"`
	Photo   string      `json:"photo,omitempty"`
}
