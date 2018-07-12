# content_type_helper
This is a pre-ware so to install you'll need to change your `app.go`:

```
app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_coke_session",
			PreWares:    []buffalo.PreWare{content_type_helper.AutoSetContentType},
		})
```

### Pre-wares:

  1. **AutoSetContentType:** This will look for file extentions on the end of a path and handle them properly.
