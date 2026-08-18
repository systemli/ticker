package storage

import "gorm.io/gorm"

// QueryOpt is a GORM query option that can be passed to store reads to apply
// preloads, ordering, or any other gorm.DB chain.
type QueryOpt = func(*gorm.DB) *gorm.DB
