package util

import (
	"regexp"
	"strconv"
)

// Validate
type Validate struct {
	i string
	r bool
	E string
}

// Validator the validator initialize func
func Validator(u string) Validate {
	var n Validate
	n.i = u
	n.r = true
	n.E = ""
	return n
}

func (s *Validate) Required() *Validate {
	if s.r == true {
		if s.i == "" || len(s.i) <= 0 {
			s.r = false
			s.E = "Is Required"
		}
	}
	return s
}

func (s *Validate) MinLength(length int) *Validate {
	if s.r == true {
		if len(s.i) < length {
			s.r = false
			s.E = "Minimum length " + strconv.Itoa(length) + " characters allowed"
		}
	}
	return s
}

func (s *Validate) MaxLength(length int) *Validate {
	if s.r == true {
		if len(s.i) > length {
			s.r = false
			s.E = "Maximum length " + strconv.Itoa(length) + " characters allowed"
		}
	}
	return s
}

func (s *Validate) IsEmail() *Validate {
	if s.r == true {
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		var match = re.MatchString(s.i)
		if !match {
			s.r = false
			s.E = "invalid email"
		}
	}
	return s
}

func (s *Validate) OneLowerCase() *Validate {
	if s.r == true {
		re := regexp.MustCompile("[a-z]+")
		var match = re.MatchString(s.i)
		if !match {
			s.r = false
			s.E = "does not contain atleast lowercase letter"
		}
	}
	return s
}

func (s *Validate) AllLowerCase() *Validate {
	if s.r == true {
		re := regexp.MustCompile("^[a-z]+$")
		var match = re.MatchString(s.i)
		if !match {
			s.r = false
			s.E = "does not contain lowercase letter"
		}
	}
	return s
}

func (s *Validate) OneUpperCase() *Validate {
	if s.r == true {
		re := regexp.MustCompile("[A-Z]+")
		var match = re.MatchString(s.i)
		if !match {
			s.r = false
			s.E = "does not contain atleast one uppercase letter"
		}
	}
	return s
}

func (s *Validate) AllUpperCase() *Validate {
	if s.r == true {
		re := regexp.MustCompile("^[A-Z]+$")
		var match = re.MatchString(s.i)
		if !match {
			s.r = false
			s.E = "all letters are not uppercase"
		}
	}
	return s
}

func (s *Validate) OneNumber() *Validate {
	if s.r == true {
		re := regexp.MustCompile("[0-9]+")
		var match = re.MatchString(s.i)
		if !match {
			s.r = false
			s.E = "does not contain atleast one numeric character"
		}
	}
	return s
}

func (s *Validate) IsSpecialCharacter() *Validate {
	if s.r == true {
		re := regexp.MustCompile("\\`|\\~|\\!|\\@|\\#|\\$|\\%|\\^|\\&|\\*|\\(|\\)|\\+|\\=|\\[|\\{|\\]|\\}|\\||\\|\\'|\\<|\\,|\\.|\\>|\\?|\\/|\"|\\;|\\:|\\s")
		var match = re.MatchString(s.i)
		if !match {
			s.r = false
			s.E = "does not contain atleast one special character"
		}
	}
	return s
}

func (s *Validate) Check() bool {
	return s.r
}
