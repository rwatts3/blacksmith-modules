module github.com/nunchistudio/blacksmith-modules/sqlike

go 1.16

require (
	github.com/flosch/pongo2/v4 v4.0.2
	github.com/nunchistudio/blacksmith v0.17.1
	github.com/segmentio/ksuid v1.0.3
	github.com/sirupsen/logrus v1.8.1
)

replace golang.org/x/net => golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20210415045647-66c3f260301c

replace golang.org/x/sync => golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9

replace github.com/nunchistudio/blacksmith => ../../blacksmith
