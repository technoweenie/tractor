module github.com/manifold/tractor

go 1.12

require (
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/d5/tengo v1.24.3
	github.com/dustin/go-jsonpointer v0.0.0
	github.com/getlantern/systray v0.0.0-20191111190243-1a6b33f30317
	github.com/gliderlabs/com v0.1.1-0.20191023181249-02615ad445ac
	github.com/gliderlabs/stdcom v0.0.0-20171109193247-64a0d4e5fd86
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d
	github.com/inconshreveable/muxado v0.0.0-20160802230925-fc182d90f26e // indirect
<<<<<<< 2bc7034a390e855eab69d0b3fc19b653716a6fe3
	github.com/lucas-clemente/quic-go v0.13.1 // indirect
=======
	github.com/lucas-clemente/quic-go v0.12.1 // indirect
>>>>>>> go.mod: updated qtalk
	github.com/manifold/qtalk v0.0.0-20191117202844-f1ce2a287d67
	github.com/mitchellh/hashstructure v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/progrium/prototypes v0.0.0-20190807232325-d9b2b4ba3a4f
	github.com/rjeczalik/notify v0.9.2
	github.com/rs/xid v1.2.1
	github.com/skratchdot/open-golang v0.0.0-20190402232053-79abb63cd66e
	github.com/spf13/afero v1.2.2
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20191117063200-497ca9f6d64f // indirect
	golang.org/x/net v0.0.0-20191118183410-d06c31c94cae // indirect
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1 // indirect
)

replace github.com/dustin/go-jsonpointer => ./vnd/github.com/dustin/go-jsonpointer
