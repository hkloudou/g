## Install

``` sh
go get github.com/hkloudou/g/db
go get github.com/hkloudou/xlib/xflag
```

## http
``` go
	app := xflag.NewApp()
	db.Use()
	xface.Config(app)
	app.Action = func(ctx *xflag.Context) error {
		db.WGet().AutoMigrate(&apiUserInfo{})
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
```