package json

import (
	stdlibJson "encoding/json"
	"testing"

	"os"

	parsecJson "github.com/prataprc/goparsec/json"
	"github.com/stretchr/testify/require"
	"github.com/vektah/goparsify"
)

func TestUnmarshal(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		result, err := Unmarshal(`true`)
		require.NoError(t, err)
		require.Equal(t, true, result)

		result, err = Unmarshal(`false`)
		require.NoError(t, err)
		require.Equal(t, false, result)

		result, err = Unmarshal(`null`)
		require.NoError(t, err)
		require.Equal(t, nil, result)

		result, err = Unmarshal(`"true"`)
		require.NoError(t, err)
		require.Equal(t, "true", result)
	})

	t.Run("array", func(t *testing.T) {
		result, err := Unmarshal(`[true, null, false, -1.23e+4]`)
		require.NoError(t, err)
		require.Equal(t, []interface{}{true, nil, false, -1.23e+4}, result)
	})

	t.Run("object", func(t *testing.T) {
		result, err := Unmarshal(`{"true":true, "false":false, "null": null, "number": 404} `)
		require.NoError(t, err)
		require.Equal(t, map[string]interface{}{"true": true, "false": false, "null": nil, "number": int64(404)}, result)
	})
}

func BenchmarkUnmarshalParsec(b *testing.B) {
	bytes := []byte(benchmarkString)

	for i := 0; i < b.N; i++ {
		scanner := parsecJson.NewJSONScanner(bytes)
		_, remaining := parsecJson.Y(scanner)

		require.True(b, remaining.Endof())
	}
}

func BenchmarkUnmarshalParsify(b *testing.B) {
	goparsify.EnableLogging(os.Stdout)
	for i := 0; i < b.N; i++ {
		_, err := Unmarshal(benchmarkString)
		require.NoError(b, err)
	}
	goparsify.DumpDebugStats()
}

func BenchmarkUnmarshalStdlib(b *testing.B) {
	bytes := []byte(benchmarkString)
	var result interface{}
	for i := 0; i < b.N; i++ {
		err := stdlibJson.Unmarshal(bytes, &result)
		require.NoError(b, err)
	}
}

// This string was taken from http://json.org/example.html
const benchmarkString = `{"web-app": {
  "servlet": [
    {
      "servlet-name": "cofaxCDS",
      "servlet-class": "org.cofax.cds.CDSServlet",
      "init-param": {
        "configGlossary:installationAt": "Philadelphia, PA",
        "configGlossary:adminEmail": "ksm@pobox.com",
        "configGlossary:poweredBy": "Cofax",
        "configGlossary:poweredByIcon": "/images/cofax.gif",
        "configGlossary:staticPath": "/content/static",
        "templateProcessorClass": "org.cofax.WysiwygTemplate",
        "templateLoaderClass": "org.cofax.FilesTemplateLoader",
        "templatePath": "templates",
        "templateOverridePath": "",
        "defaultListTemplate": "listTemplate.htm",
        "defaultFileTemplate": "articleTemplate.htm",
        "useJSP": false,
        "jspListTemplate": "listTemplate.jsp",
        "jspFileTemplate": "articleTemplate.jsp",
        "cachePackageTagsTrack": 200,
        "cachePackageTagsStore": 200,
        "cachePackageTagsRefresh": 60,
        "cacheTemplatesTrack": 100,
        "cacheTemplatesStore": 50,
        "cacheTemplatesRefresh": 15,
        "cachePagesTrack": 200,
        "cachePagesStore": 100,
        "cachePagesRefresh": 10,
        "cachePagesDirtyRead": 10,
        "searchEngineListTemplate": "forSearchEnginesList.htm",
        "searchEngineFileTemplate": "forSearchEngines.htm",
        "searchEngineRobotsDb": "WEB-INF/robots.db",
        "useDataStore": true,
        "dataStoreClass": "org.cofax.SqlDataStore",
        "redirectionClass": "org.cofax.SqlRedirection",
        "dataStoreName": "cofax",
        "dataStoreDriver": "com.microsoft.jdbc.sqlserver.SQLServerDriver",
        "dataStoreUrl": "jdbc:microsoft:sqlserver://LOCALHOST:1433;DatabaseName=goon",
        "dataStoreUser": "sa",
        "dataStorePassword": "dataStoreTestQuery",
        "dataStoreTestQuery": "SET NOCOUNT ON;select test='test';",
        "dataStoreLogFile": "/usr/local/tomcat/logs/datastore.log",
        "dataStoreInitConns": 10,
        "dataStoreMaxConns": 100,
        "dataStoreConnUsageLimit": 100,
        "dataStoreLogLevel": "debug",
        "maxUrlLength": 500}},
    {
      "servlet-name": "cofaxEmail",
      "servlet-class": "org.cofax.cds.EmailServlet",
      "init-param": {
      "mailHost": "mail1",
      "mailHostOverride": "mail2"}},
    {
      "servlet-name": "cofaxAdmin",
      "servlet-class": "org.cofax.cds.AdminServlet"},

    {
      "servlet-name": "fileServlet",
      "servlet-class": "org.cofax.cds.FileServlet"},
    {
      "servlet-name": "cofaxTools",
      "servlet-class": "org.cofax.cms.CofaxToolsServlet",
      "init-param": {
        "templatePath": "toolstemplates/",
        "log": 1,
        "logLocation": "/usr/local/tomcat/logs/CofaxTools.log",
        "logMaxSize": "",
        "dataLog": 1,
        "dataLogLocation": "/usr/local/tomcat/logs/dataLog.log",
        "dataLogMaxSize": "",
        "removePageCache": "/content/admin/remove?cache=pages&id=",
        "removeTemplateCache": "/content/admin/remove?cache=templates&id=",
        "fileTransferFolder": "/usr/local/tomcat/webapps/content/fileTransferFolder",
        "lookInContext": 1,
        "adminGroupID": 4,
        "betaServer": true}}],
  "servlet-mapping": {
    "cofaxCDS": "/",
    "cofaxEmail": "/cofaxutil/aemail/*",
    "cofaxAdmin": "/admin/*",
    "fileServlet": "/static/*",
    "cofaxTools": "/tools/*"},

  "taglib": {
    "taglib-uri": "cofax.tld",
    "taglib-location": "/WEB-INF/tlds/cofax.tld"}}}`
