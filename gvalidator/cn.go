package gvalidator

func init() {
	registe("CN-CreditCode", `^[1|5|9][1|2|3]\d{6}[^_IOZSVa-z\W]{10}$`)
	registe("CN-IdCard", `(^[1-9]\d{7}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])\d{3}$)|(^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])\d{3}(\d|X)$)`)
	registe("CN-Mobile", `^1[3-9]\d{9}$`)
}
