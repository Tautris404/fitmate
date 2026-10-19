package services

import "errors"

var (
	ErrInvalidSwiperID      = errors.New("invalid swiper id")
	ErrInvalidSwipeTargetID = errors.New("invalid swipe target id")
	ErrSwipeSelf            = errors.New("cannot swipe yourself")
	ErrSwipeActionInvalid   = errors.New("swipe action must be patinka or nepatinka")
	ErrSwipeTargetNotFound  = errors.New("athlete not found")
	ErrSwipeExists          = errors.New("athlete has already been swiped")
	ErrFailedToCreateSwipe  = errors.New("failed to create swipe")
)
