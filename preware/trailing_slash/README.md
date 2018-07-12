# Trailing Slash Helper
This is a pre-ware so to install you'll need to change your `app.go`:

```
app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_coke_session",
			PreWares:    []buffalo.PreWare{trailing_slash.RemoveTrailingSlash},
		})
```

### Pre-wares:

  1. **RemoveTrailingSlash:** This will check for if the requested path has a trailing slash and remove it.
