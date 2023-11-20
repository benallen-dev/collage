# Collage
A basic web app for collecting images during a presentation or meeting.

Users can paste images and add their name to their submissions. Presenters can see who has submitted images, and reveal them all to the session.

Presenters can also reset the session.

## Requirements

You must have [Go templ](https://templ.guide) installed. Templ is used to render jsx-like templates into code that returns HTML.

## How to run

```sh
$ templ generate
$ go run cmd/main.go
```

Then go to `localhost:1323` in your browser.

## Running in development

Install [Air - Live reload for Go apps](https://github.com/cosmtrek/air), then

```sh
$ air
```

## Simple.css

Simple.css was the first no-class drop-in stylesheet I found and it looks pretty good.

