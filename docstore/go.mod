module github.com/nunchistudio/blacksmith-modules/docstore

go 1.16

require (
	github.com/nunchistudio/blacksmith v0.17.0
	github.com/sirupsen/logrus v1.8.1
	gocloud.dev v0.22.0
	gocloud.dev/docstore/mongodocstore v0.22.0
)

replace golang.org/x/net => golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20210415045647-66c3f260301c

replace golang.org/x/sync => golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9