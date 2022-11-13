# Go Dynamic Web Sample

## Description
This sample project is a basic Go web server with server side rendered html. Dynamic components are enabled via html fragments sent over a websocket. On the client side the dom is updated via [morphdom](https://github.com/patrick-steele-idem/morphdom).

This approach is similar to [Phoenix LiveView](https://github.com/phoenixframework/phoenix_live_view), [Hotwire](https://hotwired.dev/), and others.

Tailwind CSS is used for the basic styling and esbuild for bundling.

## Running Sample

Start tailwind watcher.
```sh
npm run tailwind
```

Run esbuild to bundle javascript
```sh
npm run esbuild
```

Run go server
```sh
go run main.go
```

If you visit http://localhost:8080 you will see a timestamp that updates every second and the live number of connections. Open another instance in a different tab to see the connection count go up.
<img width="632" alt="image" src="https://user-images.githubusercontent.com/16405245/201549802-c093b3ae-a55e-4b4e-a366-a17adb49f19b.png">
