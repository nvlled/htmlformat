package htmlformat

import (
	"testing"
)

func TestCollapseLeftWhitespace(t *testing.T) {
	cases := [][]string{
		{"", ""},
		{" ", " "},
		{"\t", " "},
		{"\n", "\n"},
		{"\n\t ", "\n"},
		{"x", "x"},
		{" x", " x"},
		{"\tx", " x"},
		{"\nx", "\nx"},
		{"xy ", "xy "},
		{"  xy ", " xy "},
		{"\t xy ", " xy "},
		{"\t xy ", " xy "},
		{"  \n xy ", "\nxy "},
	}

	for _, c := range cases {
		arg, expected := c[0], c[1]
		actual := collapseLeftWhitespace(arg)
		if actual != expected {
			t.Errorf("expected=%q, got=%q", expected, actual)
		}
	}

}

func TestCollapseRightWhitespace(t *testing.T) {
	cases := [][]string{
		{"", ""},
		{" ", " "},
		{"\t", " "},
		{"\n", "\n"},
		{"\n\t ", "\n"},
		{"x", "x"},
		{" x ", " x "},
		{"x\n", "x\n"},
		{"xy ", "xy "},
		{"xy        ", "xy "},
		{"xy  \t\t  ", "xy "},
		{"\nxy\t\n\t\t  ", "\nxy\n"},
		{" xy  ", " xy "},
		{" xy\t", " xy "},
		{" xy\t ", " xy "},
		{"\nxy\n", "\nxy\n"},
	}

	for _, c := range cases {
		arg, expected := c[0], c[1]
		actual := collapseRightWhitespace(arg)
		if actual != expected {
			t.Errorf("expected=%q, got=%q", expected, actual)
		}
	}

}

func TestCollapseWhitespace(t *testing.T) {
	cases := [][]string{
		{"", ""},
		{" ", " "},
		{"\t", " "},
		{"\n", "\n"},
		{"\n\t ", "\n"},
		{"x", "x"},
		{"  x  ", " x "},
		{"\nx\n", "\nx\n"},
		{"\txy ", " xy "},
		{"\txy        ", " xy "},
		{"\txy  \t\t  ", " xy "},
		{" \t\t  xy \t\t\n", " xy\n"},
	}

	for _, c := range cases {
		arg, expected := c[0], c[1]
		actual := collapseWhitespace(arg)
		if actual != expected {
			t.Errorf("expected=%q, got=%q", expected, actual)
		}
	}

}

func TestDedent(t *testing.T) {

	println("----------")
	println(dedent(`
	function foo(x)
		if x % 2 == 0 then
			if blah then
				return 2
			else
				return 3
			end
			return 1
		end
		return 0
	end
	`))

}

func TestFormatSimple(t *testing.T) {
	/*
				<p>test</p>
				<body><div id="site-menu-container"><ul id="site-menu"><li><a class="selected"href="/">/top/</a></li><li><a class=""href="/?feed=new">/new/</a></li><li><a class=""href="/?feed=best">/best/</a></li><li><a class=""href="/?feed=ask">/ask/</a></li><li><a class=""href="/?feed=show">/show/</a></li><li><a class=""href="/?feed=job">/job/</a></li></ul></div><div id="site-nav"><div id="site-logo"><a href="/">^</a></div><a id="site-name"href="/">slacker news</a><div id="site-nav-spacing"></div><div><div id="account-info"><li><a href="#/login">login</a></li><li><a href="/about">about</a></li></div></div></div><div id="wrapper"><div>
				<script> function foo() {
					return 1+1;
				}</script>
				</div>
				</div>
				</div>

			<em>x</em>  <em>y</em>
		</p>
		<p>test 1</p>
		<p> test 2
		</p>

		<p>
		test 3
		</p>
	*/
	output := String(`
	
	<!DOCTYPE html>
<html lang="en">
<head>
    <!-- Optimizely -->
<script type="text/javascript">
    window.optimizely = window.optimizely || [];

    function checkCookieConsent() {
        // Check if the "cookie_consent" cookie exists
        if (document.cookie.indexOf('cookie_consent') === -1) {
            window.optimizely.push({type: "holdEvents"});
            return;
        }

        // Get the value of the "cookie_consent" cookie
        var cookies = document.cookie.split(';');
        var cookieVal;
        for (var i = 0; i < cookies.length; i++) {
            var cookie = cookies[i].trim();
            if (cookie.indexOf('cookie_consent=') === 0) {
                cookieVal = cookie.substring('cookie_consent='.length);
                break;
            }
        }

        // Check if the value includes 'analytics_storage'
        if (cookieVal && cookieVal.indexOf('analytics_storage') !== -1) {
            // If true, send events and stop polling
            window.optimizely.push({type: "sendEvents"});
            clearInterval(pollInterval);
        }
    }

    // Poll every 500ms to check if the "cookie_consent" cookie exists
    var pollInterval = setInterval(checkCookieConsent, 500);

    // Check on initial load
    checkCookieConsent();
</script>

<script src="https://resources.jetbrains.com/storage/optly/26613100737.js">
    // www.jetbrains.com
</script>
<!-- End Optimizely --><!-- Error reporting -->
<script>(function(){
  window.reportError = function(msg, file, line, col, err, isUnhandledRejection){};
  var prevOnError = window.onerror;
  var onError = function(msg, file, line, col, err) {
    reportError(msg, file, line, col, err, false);
    prevOnError && prevOnError.apply(window, arguments);
    return false;
  };
  window.onerror = onError;
  // Setup reporting for unhandled Promise rejection errors
  window.addEventListener("unhandledrejection", function(e) {
    if (!e.reason) return;
    var l = getSrcLocation(e.reason);
    reportError(e.reason.message, l.file, l.line, l.col, e.reason, true);
  });
  // Setup reporting for console.error and console.warn calls
  patchConsole('error');
  patchConsole('warn');
  // Utility functions
  function patchConsole(fnName) {
    var fn = console[fnName];
    console[fnName] = function() {
      fn.apply(console, arguments);
      var l; try {
        throw new Error('_');
      } catch (err) {
        l = getSrcLocation(err, 1);
      }
      var msg = 'console.' + fnName + ': ' + Array.prototype.join.call(arguments, ' ');
      reportError(msg, l.file, l.line, l.col, undefined, false);
    };
  }
  function getSrcLocation(err, sd) {
    var s = err && err.stack;
    var l = s && s.split("\n")[1 + (sd|0)];
    var r = l && (/^\s*at [^(]*\((.*?):(\d+)(:\d+)?\)$/.exec(l) || /^\s*at (.*?):(\d+)(:\d+)?$/.exec(l));
    return r ? {file: r[1], line: r[2], col: r[3]} : {};
  }
})();</script>
<!-- Error reporting --><!-- Google Tag Manager -->
<script>(function() {
  // Initialize Tag Manager queue
  window.dataLayer = window.dataLayer || [];
  window.gtmLoaded = false;
  // Setup reporting for errors that occurred before Tag Manager initialized
  var prevReportError = window.reportError;
  var reportError = function(msg, file, line, col, err, isUnhandledRejection) {
    if (!window.gtmLoaded || isUnhandledRejection) {
      // Reproduce the behavior of the Tag Manager error handler
      window.dataLayer.push(makeEvt(msg, file, line));
    }
    prevReportError && prevReportError.apply(window, arguments);
  };
  window.reportError = reportError;
  // Utility functions
  function makeEvt(msg, file, line) {
    return {
      event: "gtm.pageError", "gtm.errorMessage": msg,
      "gtm.errorUrl": file, "gtm.errorLineNumber": line
    };
  }
})();</script>
<script>(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':
  new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],
  j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
  'https://www.googletagmanager.com/gtm.js?id='+i+dl;j.addEventListener(
    'load', function(){window.gtmLoaded=true});f.parentNode.insertBefore(j,f);
})(window,document,'script','dataLayer','GTM-5P98');</script>
<!-- End Google Tag Manager -->
<title>Thank you for downloading JetBrains Rider!</title>

<meta charset="utf-8">
<meta http-equiv="x-ua-compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, maximum-scale=1">


<link rel="icon" href="/favicon.ico?r=1234" type="image/x-icon"><!-- 48×48 -->
<link rel="icon" href="/icon.svg?r=1234" type="image/svg+xml" sizes="any">
<link rel="apple-touch-icon" href="/apple-touch-icon.png?r=1234" sizes="180x180"><!-- 180×180 -->
<link rel="icon" href="/icon-512.png?r=1234" type="image/png" sizes="512x512">
<link rel="manifest" href="/site.webmanifest" crossorigin="use-credentials">

<meta name="apple-mobile-web-app-title" content="JetBrains">
<meta name="application-name" content="JetBrains">
<meta name="msapplication-TileColor" content="#000000">
<meta name="theme-color" content="#000000">
<link rel="canonical" href="https://www.jetbrains.com/rider/download/download-thanks.html"/><!-- .420-->
<meta name="description" content=".NET IDE based on the IntelliJ platform and ReSharper. Supports C#, ASP.NET, ASP.NET MVC, .NET Core, Unity and Xamarin
"/>


    <link rel="alternate" hreflang="x-default" href="https://www.jetbrains.com/rider/download/download-thanks.html" />

                        <link rel="alternate" hreflang="en" href="https://www.jetbrains.com/rider/download/download-thanks.html" />
                                                            <link rel="alternate" hreflang="en-CN" href="https://www.jetbrains.com.cn/en-us/rider/download/download-thanks.html" />
                                                                        <link rel="alternate" hreflang="de" href="https://www.jetbrains.com/de-de/rider/download/download-thanks.html" />
                                                                                                    <link rel="alternate" hreflang="es" href="https://www.jetbrains.com/es-es/rider/download/download-thanks.html" />
                                                                                                    <link rel="alternate" hreflang="fr" href="https://www.jetbrains.com/fr-fr/rider/download/download-thanks.html" />
                                                                                                    <link rel="alternate" hreflang="ja" href="https://www.jetbrains.com/ja-jp/rider/download/download-thanks.html" />
                                                                                                    <link rel="alternate" hreflang="ko" href="https://www.jetbrains.com/ko-kr/rider/download/download-thanks.html" />
                                                                                                    <link rel="alternate" hreflang="ru" href="https://www.jetbrains.com/ru-ru/rider/download/download-thanks.html" />
                                                                                                    <link rel="alternate" hreflang="zh-Hans" href="https://www.jetbrains.com/zh-cn/rider/download/download-thanks.html" />
                                            <link rel="alternate" hreflang="zh-CN" href="https://www.jetbrains.com.cn/rider/download/download-thanks.html" />
                                                                                        <link rel="alternate" hreflang="pt-BR" href="https://www.jetbrains.com/pt-br/rider/download/download-thanks.html" />
                                                                                
<meta name="robots" content="noindex"/>



<script>
    default_site_language = 'en-us';
    var current_lang = 'en-us';
                var i18n_info = {"current_lang": "en-us", "languages": [{"canonical": "en", "code": "en-us", "label": "English", "page_translated": true, "url": "/rider/download/download-thanks.html"}, {"canonical": "de", "code": "de-de", "label": "Deutsch", "page_translated": true, "url": "/de-de/rider/download/download-thanks.html"}, {"canonical": "es", "code": "es-es", "label": "Espa\u00f1ol", "page_translated": true, "url": "/es-es/rider/download/download-thanks.html"}, {"canonical": "fr", "code": "fr-fr", "label": "Fran\u00e7ais", "page_translated": true, "url": "/fr-fr/rider/download/download-thanks.html"}, {"canonical": "ja", "code": "ja-jp", "label": "\u65e5\u672c\u8a9e", "page_translated": true, "url": "/ja-jp/rider/download/download-thanks.html"}, {"canonical": "ko", "code": "ko-kr", "label": "\ud55c\uad6d\uc5b4", "page_translated": true, "url": "/ko-kr/rider/download/download-thanks.html"}, {"canonical": "ru", "code": "ru-ru", "label": "\u0420\u0443\u0441\u0441\u043a\u0438\u0439", "page_translated": true, "url": "/ru-ru/rider/download/download-thanks.html"}, {"canonical": "zh-Hans", "code": "zh-cn", "label": "\u7b80\u4f53\u4e2d\u6587", "page_translated": true, "url": "/zh-cn/rider/download/download-thanks.html"}, {"canonical": "pt-BR", "code": "pt-br", "label": "Portugu\u00eas do Brasil", "page_translated": true, "url": "/pt-br/rider/download/download-thanks.html"}]};
            var navigationMenu = {"primary": {"items": [{"title": "Developer Tools", "banners": [{"isActive": false, "title": "JetBrains IDEs", "description": "Make it happen. With code.", "logoSrc": "/img/banners-menu-main/ides.svg", "actionLabel": "Learn more", "url": "/ides/", "isUrlShouldBeLocalized": true, "bgColor": "#A5029E", "bgGradient": "linear-gradient(125deg, #4101A9 31.81%, #A5029E 71.18%, #EF3692 110.54%)", "cleaned_url": "/ides/"}, {"isActive": false, "title": "Qodana", "description": "The only code quality platform as smart as JetBrains IDEs", "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/qodana/qodana.svg", "actionLabel": "Learn more", "url": "/qodana/", "isUrlShouldBeLocalized": true, "bgColor": "#F02D8A", "bgGradient": "linear-gradient(208deg, #FB6540 0%, #F02D8A 24.83%, #2A017E 99.48%)", "cleaned_url": "/qodana/"}], "suggestions": [{"isActive": false, "url": "/products/", "isUrlShouldBeLocalized": true, "title": "Not sure which tool is best for you?", "description": "Whichever technologies you use, there's a JetBrains tool to match", "actionLabel": "Find your tool", "cleaned_url": "/products/"}], "submenu": {"layout": "auto-fill inline inline inline", "columns": [{"title": "JETBRAINS IDEs", "mobileLayout": "forceTwoColumns", "subColumns": [{"items": [{"isActive": false, "title": "All IDEs", "url": "/ides/", "isUrlShouldBeLocalized": true, "cleaned_url": "/ides/"}, {"isActive": false, "title": "Aqua", "url": "/aqua/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/aqua/aqua.svg", "cleaned_url": "/aqua/"}, {"isActive": false, "title": "CLion", "url": "/clion/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/clion/clion.svg", "cleaned_url": "/clion/"}, {"isActive": false, "title": "DataGrip", "url": "/datagrip/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/datagrip/datagrip.svg", "cleaned_url": "/datagrip/"}, {"isActive": false, "title": "DataSpell", "url": "/dataspell/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/dataspell/dataspell.svg", "cleaned_url": "/dataspell/"}, {"isActive": false, "title": "Fleet", "url": "/fleet/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/fleet/fleet.svg", "cleaned_url": "/fleet/"}, {"isActive": false, "title": "GoLand", "url": "/go/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/goland/goland.svg", "cleaned_url": "/go/"}]}, {"items": [{"isActive": false, "title": "IntelliJ&nbsp;IDEA", "url": "/idea/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/intellij-idea/intellij-idea.svg", "cleaned_url": "/idea/"}, {"isActive": false, "title": "PhpStorm", "url": "/phpstorm/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/phpstorm/phpstorm.svg", "cleaned_url": "/phpstorm/"}, {"isActive": false, "title": "PyCharm", "url": "/pycharm/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/pycharm/pycharm.svg", "cleaned_url": "/pycharm/"}, {"isActive": false, "title": "Rider", "url": "/rider/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/rider/rider.svg", "cleaned_url": "/rider/"}, {"isActive": false, "title": "RubyMine", "url": "/ruby/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/rubymine/rubymine.svg", "cleaned_url": "/ruby/"}, {"isActive": false, "title": "RustRover", "url": "/rust/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/rustrover/rustrover.svg", "cleaned_url": "/rust/"}, {"isActive": false, "title": "WebStorm", "url": "/webstorm/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/webstorm/webstorm.svg", "cleaned_url": "/webstorm/"}, {"isActive": false, "title": "Writerside", "url": "/writerside/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/writerside/writerside.svg", "cleaned_url": "/writerside/"}]}]}, {"title": "PLUGINS & SERVICES", "mobileLayout": "forceTwoColumns", "items": [{"isActive": false, "title": "All Plugins", "url": "https://plugins.jetbrains.com/", "cleaned_url": "https://plugins.jetbrains.com/"}, {"isActive": false, "title": "JetBrains AI", "url": "/ai/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/ai/ai.svg", "cleaned_url": "/ai/"}, {"isActive": false, "title": "IDE Themes", "url": "https://plugins.jetbrains.com/search?tags=Theme", "cleaned_url": "https://plugins.jetbrains.com/search?tags=Theme"}, {"isActive": false, "title": "Big Data Tools", "url": "https://plugins.jetbrains.com/plugin/12494-big-data-tools", "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/big-data-tools/big-data-tools.svg", "cleaned_url": "https://plugins.jetbrains.com/plugin/12494-big-data-tools"}, {"isActive": false, "title": "Code With Me", "url": "/code-with-me/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/code-with-me/code-with-me.svg", "cleaned_url": "/code-with-me/"}, {"isActive": false, "title": "RiderFlow", "url": "/riderflow/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/riderflow/riderflow.svg", "cleaned_url": "/riderflow/"}, {"isActive": false, "title": "Scala", "url": "https://plugins.jetbrains.com/plugin/1347-scala", "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/scala/scala.svg", "cleaned_url": "https://plugins.jetbrains.com/plugin/1347-scala"}, {"isActive": false, "title": "Toolbox App", "url": "/toolbox-app/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/toolbox/toolbox.svg", "cleaned_url": "/toolbox-app/"}, {"isActive": false, "title": "Grazie", "url": "/grazie/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/grazie/grazie.svg", "cleaned_url": "/grazie/"}]}, {"title": ".NET & VISUAL STUDIO", "hasSeparator": true, "items": [{"isActive": false, "title": "Rider", "url": "/rider/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/rider/rider.svg", "cleaned_url": "/rider/"}, {"isActive": false, "title": "ReSharper", "url": "/resharper/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/resharper/resharper.svg", "cleaned_url": "/resharper/"}, {"isActive": false, "title": "ReSharper C++", "url": "/resharper-cpp/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/resharper-cpp/resharper-cpp.svg", "cleaned_url": "/resharper-cpp/"}, {"isActive": false, "title": "dotCover", "url": "/dotcover/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/dotcover/dotcover.svg", "cleaned_url": "/dotcover/"}, {"isActive": false, "title": "dotMemory", "url": "/dotmemory/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/dotmemory/dotmemory.svg", "cleaned_url": "/dotmemory/"}, {"isActive": false, "title": "dotPeek", "url": "/decompiler/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/dotpeek/dotpeek.svg", "cleaned_url": "/decompiler/"}, {"isActive": false, "title": "dotTrace", "url": "/profiler/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/dottrace/dottrace.svg", "cleaned_url": "/profiler/"}, {"isActive": false, "title": ".NET Tools Plugins", "url": "https://plugins.jetbrains.com/search?isFeaturedSearch=true&products=resharper&products=rider", "cleaned_url": "https://plugins.jetbrains.com/search?isFeaturedSearch=true&products=resharper&products=rider"}]}, {"title": "LANGUAGES & FRAMEWORKS", "hasSeparator": true, "items": [{"isActive": false, "title": "Kotlin", "url": "https://kotlinlang.org/", "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/kotlin/kotlin.svg", "cleaned_url": "https://kotlinlang.org/"}, {"isActive": false, "title": "Ktor", "url": "https://ktor.io/", "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/ktor/ktor.svg", "cleaned_url": "https://ktor.io/"}, {"isActive": false, "title": "MPS", "url": "/mps/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/mps/mps.svg", "cleaned_url": "/mps/"}, {"isActive": false, "title": "Compose Multiplatform", "url": "/compose-multiplatform/", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/compose-multiplatform/compose-multiplatform.svg", "cleaned_url": "/compose-multiplatform/"}]}]}, "priority": 3}, {"title": "Team Tools", "banners": [{"isActive": false, "title": "Datalore", "description": "A collaborative data science platform. Available online and on-premises", "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/datalore/datalore.svg", "actionLabel": "Learn more", "url": "/datalore/", "isUrlShouldBeLocalized": true, "bgColor": "#005CD1", "bgGradient": "linear-gradient(120.81deg, #003396 11.31%, #009CF4 95.37%)", "cleaned_url": "/datalore/"}, {"isActive": false, "title": "YouTrack", "description": "Powerful project management for all your teams", "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/youtrack/youtrack.svg", "actionLabel": "Learn more", "url": "/youtrack/", "isUrlShouldBeLocalized": true, "bgColor": "#6B57FF", "bgGradient": "linear-gradient(313deg, #D919D0 10.26%, #BC003C 91.89%)", "cleaned_url": "/youtrack/"}], "submenu": {"layout": "8 4", "columns": [{"title": "IN-CLOUD AND ON-PREMISES SOLUTIONS", "subColumns": [{"items": [{"isActive": false, "title": "Datalore", "url": "/datalore/", "isUrlShouldBeLocalized": true, "description": "A collaborative data science platform", "cleaned_url": "/datalore/"}, {"isActive": false, "title": "TeamCity", "url": "/teamcity/", "isUrlShouldBeLocalized": true, "description": "Powerful Continuous Integration out of the box", "cleaned_url": "/teamcity/"}, {"isActive": false, "title": "CodeCanvas", "url": "/codecanvas/", "isUrlShouldBeLocalized": true, "description": "Cloud development environments for your infrastructure", "cleaned_url": "/codecanvas/"}]}, {"items": [{"isActive": false, "title": "YouTrack", "url": "/youtrack/", "isUrlShouldBeLocalized": true, "description": "Powerful project management for all your teams", "cleaned_url": "/youtrack/"}, {"isActive": false, "title": "Qodana", "url": "/qodana/", "isUrlShouldBeLocalized": true, "description": "The code quality platform for teams", "cleaned_url": "/qodana/"}]}]}, {"title": "EXTENSIONS", "hasSeparator": true, "items": [{"isActive": false, "title": "TeamCity Plugins", "url": "https://plugins.jetbrains.com/teamcity/", "cleaned_url": "https://plugins.jetbrains.com/teamcity/"}, {"isActive": false, "title": "YouTrack Extensions", "url": "https://plugins.jetbrains.com/youtrack/", "cleaned_url": "https://plugins.jetbrains.com/youtrack/"}, {"isActive": false, "title": "JetBrains Hub", "url": "/hub/", "isUrlShouldBeLocalized": true, "cleaned_url": "/hub/"}]}]}, "priority": 2}, {"title": "Education", "banners": [{"isActive": false, "title": "JetBrains Academy", "description": "Find your way in learning or teaching computer science", "actionLabel": "Discover more", "url": "/academy", "isUrlShouldBeLocalized": true, "logoSrc": "${RESOURCES_URL_PLACEHOLDER}/storage/logos/web/jetbrains-academy/jetbrains-academy.svg", "bgColor": "#B01DF6", "bgGradient": "linear-gradient(335.07deg, #636CEA 0%, #834CEF 40.63%, #771F89 100%)", "cleaned_url": "/academy"}], "submenu": {"columns": [{"title": "FOR LEARNERS", "layout": "11 11 11", "subColumns": [{"items": [{"isActive": false, "title": "Programming languages", "url": "/academy/", "isUrlShouldBeLocalized": true, "description": "Select a language and try different approaches to learning it", "cleaned_url": "/academy/"}, {"isActive": false, "title": "University relations", "url": "/education/university-relations/", "isUrlShouldBeLocalized": true, "description": "Study offline with academic programs", "cleaned_url": "/education/university-relations/"}, {"isActive": false, "title": "Internships", "url": "/careers/internships/", "isUrlShouldBeLocalized": true, "description": "Apply for internships and flexible jobs for students\n", "cleaned_url": "/careers/internships/"}]}]}, {"title": "FOR EDUCATORS", "layout": "11 11 11", "subColumns": [{"items": [{"isActive": false, "title": "Teaching with JetBrains IDEs", "url": "/academy/teaching/", "isUrlShouldBeLocalized": true, "description": "Create courses and share your knowledge", "cleaned_url": "/academy/teaching/"}, {"isActive": false, "title": "Kotlin for education", "url": "https://kotlinlang.org/education/", "isUrlShouldBeLocalized": true, "description": "Teach a wide range of Kotlin courses", "cleaned_url": "https://kotlinlang.org/education/"}]}, {"title": "FOR TEAMS", "items": [{"isActive": false, "title": "Professional development", "url": "https://lp.jetbrains.com/academy/for-organizations", "isUrlShouldBeLocalized": true, "description": "Ensure your team has up-to-date technical skills", "cleaned_url": "https://lp.jetbrains.com/academy/for-organizations"}]}]}, {"title": "FREE LICENSES", "hasSeparator": true, "items": [{"isActive": false, "title": "For students and teachers", "url": "/community/education/#students/", "isUrlShouldBeLocalized": true, "description": "JetBrains IDEs for individual academic use", "cleaned_url": "/community/education/#students/"}, {"isActive": false, "title": "For educational institutions", "url": "/community/education/#classrooms", "isUrlShouldBeLocalized": true, "description": "JetBrains IDEs and team tools for classroom use", "cleaned_url": "/community/education/#classrooms"}, {"isActive": false, "title": "For bootcamps and courses", "url": "/academy/bootcamps/", "isUrlShouldBeLocalized": true, "description": "JetBrains IDEs for your students", "cleaned_url": "/academy/bootcamps/"}]}]}, "priority": 1}, {"title": "Solutions", "banners": [{"isActive": false, "title": "Developer Tools for Your Business", "description": "Professional tools for productive development", "actionLabel": "Learn more", "url": "/business/", "isUrlShouldBeLocalized": true, "logoSrc": "/img/banners-menu-main/containers.svg", "bgColor": "#6B57FF", "bgGradient": "linear-gradient(246.1deg, rgb(0 224 214) 1.67%, rgb(126 27 253) 92.48%)", "cleaned_url": "/business/"}, {"isActive": false, "title": "Remote Development", "description": "Connect to remote dev environments from anywhere in seconds", "actionLabel": "Discover more", "url": "/remote-development/", "isUrlShouldBeLocalized": true, "bgColor": "#2DF388", "bgGradient": "linear-gradient(240.88deg, #2DF388 0%, #05BF87 37.75%, #027474 98.39%)", "cleaned_url": "/remote-development/"}], "submenu": {"layout": "8 4", "columns": [{"title": "BY INDUSTRY & TECHNOLOGY", "layout": "6 6", "subColumns": [{"items": [{"isActive": false, "title": "Remote Development", "url": "/remote-development/", "isUrlShouldBeLocalized": true, "description": "Tools for remote development for you and your team", "cleaned_url": "/remote-development/"}, {"isActive": false, "title": "Game Development", "url": "/gamedev/", "isUrlShouldBeLocalized": true, "description": "Tools for game development for any platform", "cleaned_url": "/gamedev/"}, {"isActive": false, "title": "DevOps", "url": "/devops/", "isUrlShouldBeLocalized": true, "description": "Tools and integrations for any infrastructure", "cleaned_url": "/devops/"}, {"isActive": false, "title": "Multiplatform Development", "url": "/kotlin-multiplatform/", "isUrlShouldBeLocalized": true, "description": "Flexible cross-platform development with Kotlin", "cleaned_url": "/kotlin-multiplatform/"}]}, {"items": [{"isActive": false, "title": "AI Service and AI Assistant", "url": "/ai/", "isUrlShouldBeLocalized": true, "description": "Augmented developer environments and team tools", "cleaned_url": "/ai/"}, {"isActive": false, "title": "C++ Tools", "url": "/cpp/", "isUrlShouldBeLocalized": true, "description": "Tools for C/C++ development for any platform", "cleaned_url": "/cpp/"}, {"isActive": false, "title": "Data Tools", "url": "/data-tools/", "isUrlShouldBeLocalized": true, "description": "Tools for Big Data and Data Science", "cleaned_url": "/data-tools/"}, {"isActive": false, "title": "License Vault", "url": "/license-vault/", "isUrlShouldBeLocalized": true, "description": "Efficient management of JetBrains licenses", "cleaned_url": "/license-vault/"}, {"isActive": false, "title": "JetBrains IDE Services", "url": "/ide-services/", "isUrlShouldBeLocalized": true, "description": "Developer productivity at the scale of an organization", "cleaned_url": "/ide-services/"}]}]}, {"title": "RECOMMENDED", "hasSeparator": true, "items": [{"isActive": false, "title": "JetBrains Tools for Business", "url": "/business/", "isUrlShouldBeLocalized": true, "cleaned_url": "/business/"}, {"isActive": false, "title": "All Products Pack", "url": "/all/", "isUrlShouldBeLocalized": true, "cleaned_url": "/all/"}, {"isActive": false, "title": ".NET Tools", "url": "/dotnet/", "isUrlShouldBeLocalized": true, "cleaned_url": "/dotnet/"}, {"isActive": false, "title": "All JetBrains Products", "url": "/products/", "isUrlShouldBeLocalized": true, "cleaned_url": "/products/"}, {"isActive": false, "title": "JetBrains Marketplace", "url": "https://plugins.jetbrains.com/", "cleaned_url": "https://plugins.jetbrains.com/"}]}]}, "priority": 0}, {"title": "Support", "banners": [{"isActive": false, "title": "Download and Install", "actionLabel": "Download and Install", "url": "/products/", "isUrlShouldBeLocalized": true, "logoSrc": "/img/banners-menu-main/download.svg", "bgColor": "#6B57FF", "bgGradient": "linear-gradient(294.91deg, #FF318C -50.1%, #6B57FF 97.43%)", "cleaned_url": "/products/"}, {"isActive": false, "title": "Contact us", "actionLabel": "Contact us", "url": "/company/contacts/", "isUrlShouldBeLocalized": true, "logoSrc": "/img/banners-menu-main/test-review.svg", "bgColor": "#21D789", "bgGradient": "linear-gradient(283.8deg, #087CFA 5.73%, #21D789 100%)", "cleaned_url": "/company/contacts/"}], "submenu": {"columns": [{"title": "PRODUCT & TECHNICAL SUPPORT", "layout": "12", "subColumns": [{"items": [{"isActive": false, "title": "Support Center", "url": "/support/", "isUrlShouldBeLocalized": true, "cleaned_url": "/support/"}, {"isActive": false, "title": "Product-Specific Information", "url": "/business/documents/", "isUrlShouldBeLocalized": true, "cleaned_url": "/business/documents/"}, {"isActive": false, "title": "Product Documentation", "url": "/help/", "isUrlShouldBeLocalized": true, "cleaned_url": "/help/"}, {"isActive": false, "title": "Livestreams", "url": "/company/livestreams/", "isUrlShouldBeLocalized": true, "cleaned_url": "/company/livestreams/"}, {"isActive": false, "title": "Newsletters", "url": "/resources/newsletters/", "isUrlShouldBeLocalized": true, "cleaned_url": "/resources/newsletters/"}, {"isActive": false, "title": "Early Access", "url": "/resources/eap/", "isUrlShouldBeLocalized": true, "cleaned_url": "/resources/eap/"}, {"isActive": false, "title": "Blog", "url": "https://blog.jetbrains.com/", "isUrlShouldBeLocalized": true, "cleaned_url": "https://blog.jetbrains.com/"}]}]}, {"title": "FREQUENT TASKS", "hasSeparator": true, "items": [{"isActive": false, "title": "Manage your account", "url": "https://account.jetbrains.com/profile-details", "cleaned_url": "https://account.jetbrains.com/profile-details"}, {"isActive": false, "title": "Manage your licenses", "url": "https://account.jetbrains.com/licenses", "cleaned_url": "https://account.jetbrains.com/licenses"}, {"isActive": false, "title": "Contact Sales", "url": "/support/sales/", "isUrlShouldBeLocalized": true, "cleaned_url": "/support/sales/"}, {"isActive": false, "title": "Licensing FAQ", "url": "https://sales.jetbrains.com", "isUrlShouldBeLocalized": true, "cleaned_url": "https://sales.jetbrains.com"}]}]}, "priority": 2}, {"title": "Store", "banners": [{"isActive": false, "title": "All Products Pack", "description": "Get all JetBrains desktop tools including 10&nbsp;IDEs,<br />2&nbsp;profilers, and 3&nbsp;extensions", "actionLabel": "Learn more", "url": "/all/", "isUrlShouldBeLocalized": true, "logoSrc": "/img/banners-menu-main/discount.svg", "bgColor": "#FF318C", "bgGradient": "linear-gradient(293.2deg, rgb(253 13 122) 13.45%, rgb(252 100 67) 73.57%, rgb(248 158 7) 100%)", "cleaned_url": "/all/"}], "submenu": {"columns": [{"title": "DEVELOPER TOOLS", "layout": "12 12 12", "subColumns": [{"items": [{"isActive": false, "title": "For Individual Use", "url": "/store/#personal", "isUrlShouldBeLocalized": true, "cleaned_url": "/store/#personal"}, {"isActive": false, "title": "For Teams and Organizations", "url": "/store/#commercial", "isUrlShouldBeLocalized": true, "cleaned_url": "/store/#commercial"}, {"isActive": false, "title": "Special offers & programs", "url": "/store/#discounts", "isUrlShouldBeLocalized": true, "cleaned_url": "/store/#discounts"}]}, {"title": "SERVICES & PLUGINS", "items": [{"isActive": false, "title": "JetBrains AI", "url": "/ai/", "isUrlShouldBeLocalized": true, "cleaned_url": "/ai/"}, {"isActive": false, "title": "Marketplace", "url": "/store/plugins/", "isUrlShouldBeLocalized": true, "cleaned_url": "/store/plugins/"}]}, {"title": "LEARNING TOOLS", "items": [{"isActive": false, "title": "JetBrains Academy", "url": "/academy/buy/", "isUrlShouldBeLocalized": true, "cleaned_url": "/academy/buy/"}]}]}, {"title": "TEAM TOOLS", "layout": "12 12 12", "subColumns": [{"items": [{"isActive": false, "title": "TeamCity", "url": "/store/teamware#teamcity-store-section", "isUrlShouldBeLocalized": true, "cleaned_url": "/store/teamware#teamcity-store-section"}, {"isActive": false, "title": "YouTrack", "url": "/store/teamware#youtrack-store-section", "isUrlShouldBeLocalized": true, "cleaned_url": "/store/teamware#youtrack-store-section"}, {"isActive": false, "title": "Datalore", "url": "/datalore/", "isUrlShouldBeLocalized": true, "cleaned_url": "/datalore/"}, {"isActive": false, "title": "Qodana", "url": "/qodana/buy/", "isUrlShouldBeLocalized": true, "cleaned_url": "/qodana/buy/"}]}, {"title": "COLLABORATIVE DEVELOPMENT", "items": [{"isActive": false, "title": "Code With Me", "url": "/code-with-me/buy/", "isUrlShouldBeLocalized": true, "cleaned_url": "/code-with-me/buy/"}]}]}, {"title": "SALES SUPPORT", "hasSeparator": true, "items": [{"isActive": false, "title": "Contact Sales", "url": "/support/sales/", "isUrlShouldBeLocalized": true, "cleaned_url": "/support/sales/"}, {"isActive": false, "title": "Purchase Terms", "url": "/legal/docs/store/terms/", "isUrlShouldBeLocalized": true, "cleaned_url": "/legal/docs/store/terms/"}, {"isActive": false, "title": "FAQ", "url": "https://sales.jetbrains.com/", "isUrlShouldBeLocalized": true, "cleaned_url": "https://sales.jetbrains.com/"}, {"isActive": false, "title": "Partners and Resellers", "url": "/company/partners/", "isUrlShouldBeLocalized": true, "cleaned_url": "/company/partners/"}]}]}, "priority": 3}, {"isActive": false, "title": "Login", "url": "https://account.jetbrains.com/", "isMobileOnly": true, "cleaned_url": "https://account.jetbrains.com/"}]}, "secondary": {"isActive": true, "id": "product_rider", "url": "/rider/", "title": "Rider", "isThemeDark": true, "promoTag": {"title": "JetBrains IDEs", "url": "/ides/"}, "items": [{"isActive": false, "title": "Coming in new version", "url": "/rider/nextversion/", "isEapItem": true, "showEapVersion": true, "isItemHidden": true, "eapProductId": "rider", "cleaned_url": "/rider/nextversion/"}, {"isActive": false, "title": "What's New", "version": "2024.2", "url": "/rider/whatsnew/", "items": [{"isActive": false, "title": "What's New 2024.2", "version": "2024.2", "url": "/rider/whatsnew/", "cleaned_url": "/rider/whatsnew/"}, {"isActive": false, "title": "What's New 2024.1", "version": "2024.1", "url": "/rider/whatsnew/2024-1/", "cleaned_url": "/rider/whatsnew/2024-1/"}, {"isActive": false, "title": "What's New 2023.3", "version": "2023.3", "url": "/rider/whatsnew/2023-3/", "cleaned_url": "/rider/whatsnew/2023-3/"}], "cleaned_url": "/rider/whatsnew/"}, {"isActive": false, "title": "Features", "url": "/rider/features/", "cleaned_url": "/rider/features/"}, {"isActive": false, "title": "Learn", "url": "/rider/documentation/", "cleaned_url": "/rider/documentation/"}, {"isActive": false, "title": "Blog & Social", "url": "/rider/social/", "cleaned_url": "/rider/social/"}, {"isActive": false, "title": "Pricing", "url": "/rider/buy/", "type": "outlineButton", "cleaned_url": "/rider/buy/"}, {"isActive": false, "title": "Download", "url": "/rider/download/", "type": "button", "cleaned_url": "/rider/download/"}], "cleaned_url": "/rider/"}};
    
    var is_layout_adaptive = false;
        is_layout_adaptive = true;
    
    var disable_language_picker = false;
    
    var localized_domains = [{"defaultLanguage": "en", "domain": "blog.jetbrains.com", "locales": {"de-de": "de", "en-us": "en", "es-es": "es", "fr-fr": "fr", "ja-jp": "ja", "ko-kr": "ko", "pt-br": "pt-br", "ru-ru": "ru", "zh-cn": "zh-hans"}, "pathsLocalization": false, "suffixDefault": false}, {"defaultLanguage": "en-us", "domain": "lp.jetbrains.com", "locales": {"de-de": "de-de", "en-us": "en-us", "es-es": "es-es", "fr-fr": "fr-fr", "ja-jp": "ja-jp", "ko-kr": "ko-kr", "pt-br": "pt-br", "ru-ru": "ru-ru", "zh-cn": "zh-cn"}, "pathsLocalization": true, "suffixDefault": false}, {"defaultLanguage": "en-us", "domain": "sales.jetbrains.com", "locales": {"de-de": "de", "en-us": "en-gb", "es-es": "es", "fr-fr": "fr", "ja-jp": "ja", "ko-kr": "ko", "pt-br": "pt-br", "ru-ru": "ru", "zh-cn": "zh-cn"}, "pathsLocalization": true, "prefixPath": "hc", "suffixDefault": true}];

    var english_only_url_prefixes = [];
    
    var is_landing_view = false;
    
    var theme = 'light';
</script>

<script></script>



        <link href="/_assets/common.0a785d87e3c4f1289eaf.css" rel="stylesheet" type="text/css">
<link href="/_assets/default-page.6d438b89373c78cd01f3.css" rel="stylesheet" type="text/css">
<script src="/_assets/common.d065dd76ab685e44a4df.js" type="text/javascript"></script>
<script src="/_assets/default-page.5a82bea4fda23ba7c9cf.js" type="text/javascript"></script>
<script src="/_assets/rider/download/download-thanks.entry.201440311c9caca09fac.js" type="text/javascript"></script>

                        


<!-- Social Media tag Starts -->
    <!-- Open Graph data -->
    <meta property="og:title" content="Thank you for downloading JetBrains Rider!"/>

    <meta property="og:description" content=".NET IDE based on the IntelliJ platform and ReSharper. Supports C#, ASP.NET, ASP.NET MVC, .NET Core, Unity and Xamarin
"/>
    <meta property="og:image" content="https://resources.jetbrains.com/storage/products/jetbrains/img/meta/preview.png"/>

    <meta property="og:site_name" content="JetBrains"/>
    <meta property="og:type" content="website"/>
    <meta property="og:url" content="https://www.jetbrains.com/rider/download/download-thanks.html"/>
    <!-- OpenGraph End -->

    <!-- Schema.org data -->
                        <script type="application/ld+json">
                    {"@context": "http://schema.org", "@type": "SoftwareApplication", "applicationCategory": "DeveloperApplication", "author": {"@type": "Organization", "name": "JetBrains"}, "datePublished": "2021-05-26T09:00", "downloadUrl": "https://www.jetbrains.com/rider/download/", "image": "https://resources.jetbrains.com/storage/products/rider/img/meta/preview.png", "name": "Rider", "operatingSystem": "Windows, macOS, Linux", "publisher": {"@type": "Organization", "name": "JetBrains"}, "requirements": "RAM: 2 GB, Free space: 1.5 GB, Screen Resolution: 1024x768", "screenshot": "http://www.jetbrains.com/rider/img/screenshots/rider_navigation_preview@2x.png", "softwareVersion": "2021.1", "url": "https://www.jetbrains.com/rider/"}
                </script>
                    <script type="application/ld+json">
            {
                "@context": "http://schema.org",
                "@type": "WebPage",
                "@id": "https://www.jetbrains.com/rider/download/download-thanks.html#webpage",
                "url": "https://www.jetbrains.com/rider/download/download-thanks.html",
                "name": "Thank you for downloading JetBrains Rider!",
                "description": ".NET IDE based on the IntelliJ platform and ReSharper. Supports C#, ASP.NET, ASP.NET MVC, .NET Core, Unity and Xamarin",
                "image": "https://resources.jetbrains.com/storage/products/jetbrains/img/meta/jetbrains_1280x800.png"
            }</script>
    <!-- Schema.org end -->

    <!-- Twitter data -->
                                <meta name="twitter:description" content=".NET IDE based on the IntelliJ platform and ReSharper. Supports C#, ASP.NET, ASP.NET MVC, .NET Core, Unity and Xamarin
">
        
                    <meta name="twitter:title" content="Thank you for downloading JetBrains Rider!">
        
        
                    <meta name="twitter:card" content="summary_large_image">
                    <meta name="twitter:site" content="@jetbrainsrider">
                    <meta name="twitter:creator" content="@jetbrainsrider">
                    <meta name="twitter:image:src" content="https://resources.jetbrains.com/storage/products/rider/img/meta/preview.png">
                    <meta name="twitter:label1" content="Platforms:">
                    <meta name="twitter:data1" content="Windows, macOS, Linux">
                <!-- Twitter End -->
<!-- Social Media tag Ends -->
            </head>

<body class="nojs  body-adaptive page-color-lilac-purple wt-primary-map">

        <!-- Google Tag Manager (noscript) -->
<noscript><iframe src="https://www.googletagmanager.com/ns.html?id=GTM-5P98" height="0" width="0" style="display:none;visibility:hidden"></iframe></noscript>
<!-- End Google Tag Manager (noscript) -->
<script>
/*! modernizr 3.2.0 (Custom Build) | MIT *
 * http://modernizr.com/download/?-flexbox-flexboxtweener !*/
!function(e,n,t){function r(e,n){return typeof e===n}function o(){var e,n,t,o,i,s,l;for(var f in v)if(v.hasOwnProperty(f)){if(e=[],n=v[f],n.name&&(e.push(n.name.toLowerCase()),n.options&&n.options.aliases&&n.options.aliases.length))for(t=0;t<n.options.aliases.length;t++)e.push(n.options.aliases[t].toLowerCase());for(o=r(n.fn,"function")?n.fn():n.fn,i=0;i<e.length;i++)s=e[i],l=s.split("."),1===l.length?Modernizr[l[0]]=o:(!Modernizr[l[0]]||Modernizr[l[0]]instanceof Boolean||(Modernizr[l[0]]=new Boolean(Modernizr[l[0]])),Modernizr[l[0]][l[1]]=o),C.push((o?"":"no-")+l.join("-"))}}function i(e,n){return!!~(""+e).indexOf(n)}function s(e){return e.replace(/([a-z])-([a-z])/g,function(e,n,t){return n+t.toUpperCase()}).replace(/^-/,"")}function l(e,n){return function(){return e.apply(n,arguments)}}function f(e,n,t){var o;for(var i in e)if(e[i]in n)return t===!1?e[i]:(o=n[e[i]],r(o,"function")?l(o,t||n):o);return!1}function a(e){return e.replace(/([A-Z])/g,function(e,n){return"-"+n.toLowerCase()}).replace(/^ms-/,"-ms-")}function u(){return"function"!=typeof n.createElement?n.createElement(arguments[0]):b?n.createElementNS.call(n,"http://www.w3.org/2000/svg",arguments[0]):n.createElement.apply(n,arguments)}function d(){var e=n.body;return e||(e=u(b?"svg":"body"),e.fake=!0),e}function p(e,t,r,o){var i,s,l,f,a="modernizr",p=u("div"),c=d();if(parseInt(r,10))for(;r--;)l=u("div"),l.id=o?o[r]:a+(r+1),p.appendChild(l);return i=u("style"),i.type="text/css",i.id="s"+a,(c.fake?c:p).appendChild(i),c.appendChild(p),i.styleSheet?i.styleSheet.cssText=e:i.appendChild(n.createTextNode(e)),p.id=a,c.fake&&(c.style.background="",c.style.overflow="hidden",f=_.style.overflow,_.style.overflow="hidden",_.appendChild(c)),s=t(p,e),c.fake?(c.parentNode.removeChild(c),_.style.overflow=f,_.offsetHeight):p.parentNode.removeChild(p),!!s}function c(n,r){var o=n.length;if("CSS"in e&&"supports"in e.CSS){for(;o--;)if(e.CSS.supports(a(n[o]),r))return!0;return!1}if("CSSSupportsRule"in e){for(var i=[];o--;)i.push("("+a(n[o])+":"+r+")");return i=i.join(" or "),p("@supports ("+i+") { #modernizr { position: absolute; } }",function(e){return"absolute"==getComputedStyle(e,null).position})}return t}function m(e,n,o,l){function f(){d&&(delete E.style,delete E.modElem)}if(l=r(l,"undefined")?!1:l,!r(o,"undefined")){var a=c(e,o);if(!r(a,"undefined"))return a}for(var d,p,m,h,y,v=["modernizr","tspan"];!E.style;)d=!0,E.modElem=u(v.shift()),E.style=E.modElem.style;for(m=e.length,p=0;m>p;p++)if(h=e[p],y=E.style[h],i(h,"-")&&(h=s(h)),E.style[h]!==t){if(l||r(o,"undefined"))return f(),"pfx"==n?h:!0;try{E.style[h]=o}catch(g){}if(E.style[h]!=y)return f(),"pfx"==n?h:!0}return f(),!1}function h(e,n,t,o,i){var s=e.charAt(0).toUpperCase()+e.slice(1),l=(e+" "+x.join(s+" ")+s).split(" ");return r(n,"string")||r(n,"undefined")?m(l,n,o,i):(l=(e+" "+S.join(s+" ")+s).split(" "),f(l,n,t))}function y(e,n,r){return h(e,t,t,n,r)}var v=[],g={_version:"3.2.0",_config:{classPrefix:"",enableClasses:!0,enableJSClass:!0,usePrefixes:!0},_q:[],on:function(e,n){var t=this;setTimeout(function(){n(t[e])},0)},addTest:function(e,n,t){v.push({name:e,fn:n,options:t})},addAsyncTest:function(e){v.push({name:null,fn:e})}},Modernizr=function(){};Modernizr.prototype=g,Modernizr=new Modernizr;var C=[],w="Moz O ms Webkit",x=g._config.usePrefixes?w.split(" "):[];g._cssomPrefixes=x;var S=g._config.usePrefixes?w.toLowerCase().split(" "):[];g._domPrefixes=S;var _=n.documentElement,b="svg"===_.nodeName.toLowerCase(),z={elem:u("modernizr")};Modernizr._q.push(function(){delete z.elem});var E={style:z.elem.style};Modernizr._q.unshift(function(){delete E.style}),g.testAllProps=h,g.testAllProps=y,Modernizr.addTest("flexbox",y("flexBasis","1px",!0)),Modernizr.addTest("flexboxtweener",y("flexAlign","end",!0)),o(),delete g.addTest,delete g.addAsyncTest;for(var P=0;P<Modernizr._q.length;P++)Modernizr._q[P]();e.Modernizr=Modernizr}(window,document);

if (!Modernizr.flexbox && !Modernizr.flexboxtweener) {

  var $body = $('body');

  var nodesClasses = {
    wrapper: 'not-supported-browser',
    container: 'not-supported-browser__container',

    title: 'not-supported-browser__title',
    content: 'not-supported-browser__content',
    logo: 'not-supported-browser__logo'
  };

  var nodes = {
    wrapper: $('<div class="' + nodesClasses.wrapper + '"></div>'),
    title: $('<div class="' + nodesClasses.title + '">Sorry, your browser is not fully supported</div>'),
    content: $('<div class="' + nodesClasses.content + '">There may be some issues with pages layout in your current browser.<br/>Please use an alternate browser until we resolve the issues.<br/>Thank you.</div>'),
    container: $('<div class="' + nodesClasses.container + '"></div>'),
    logo: $('<div class="' + nodesClasses.logo + '"><svg class="sprite-img _jetbrains" xmlns:xlink="http://www.w3.org/1999/xlink"><use xlink:href="#jetbrains"></use></svg></div>')
  };

  $body.addClass('overflow-hidden');


  nodes.content
    .prepend(nodes.title)
    .prepend(nodes.logo);

  nodes.container
    .append(nodes.content);

  nodes.wrapper
    .append(nodes.container)
    .appendTo($body);
}
</script>
<div class="page">
    
            <div class="page__header ">
            <div class="page__header-language-suggestion" id="language-suggest-bar"></div>
<div class="page__header-country-suggestion" id="country-suggest-bar"></div>


<div class="site-header-container" id="js-site-header-container">
    <div class="site-header-stub site-header-stub--adaptive">
        <div class="wt-container site-header-stub__inner">
            <div class="site-header-stub__menu-main-skeleton-text" style="min-width: 69px"></div>
            <div class="site-header-stub__menu-main-skeleton-text" style="min-width: 81px"></div>
            <div class="site-header-stub__menu-main-skeleton-text" style="min-width: 46px"></div>
            <div class="site-header-stub__menu-main-skeleton-button"></div>
            <div class="site-header-stub__menu-main-skeleton-button"></div>
        </div>
    </div>
</div>



  <div class="menu-second _theme-dark" id="js-menu-second">
    <div class="wt-container">
        <div id="js-menu-second-mobile-wrapper" class="menu-second-mobile-wrapper wt-display-none">
            <div id="js-menu-second-mobile">
                <div class="menu-second-mobile wt-row wt-row_size_m wt-row_align-items_center wt-row_justify_between _theme-dark">
                    <div class="wt-col-inline menu-second-skeleton-text-2" style="width: 120px"></div>
                    <div class="wt-col-inline menu-second-skeleton-button" style="width: 80px"></div>
                </div>
            </div>
        </div>

        <div id="js-menu-second-desktop"
             class="menu-second-desktop wt-row wt-row_size_s wt-row_align-items_center wt-row_justify_between wt-row_nowrap wt-row-sm_wrap">
            <div class="wt-col-inline menu-second-title-box-wrapper">
                <a class="menu-second-title-box"
                   href="/rider/">
                    
                    <span class="menu-second-title-box__title wt-h3  wt-h3_theme_dark ">Rider</span>
                </a>
            </div>

            <div class="wt-col-auto-fill">
                <div class="wt-row wt-row_justify_end wt-row_align-items_center wt-row_size_0">
                                            
                                                                    
                                                                                <div class="wt-col-inline menu-item menu-second-skeleton-text-2 menu-second__link"></div>
                                                                    
                                                                                <div class="wt-col-inline menu-item menu-second-skeleton-text-2 menu-second__link"></div>
                                                                    
                                                                                <div class="wt-col-inline menu-item menu-second-skeleton-text-2 menu-second__link"></div>
                                                                    
                                                                                <div class="wt-col-inline menu-item menu-second-skeleton-text-2 menu-second__link"></div>
                                                                    
                                                                                <div class="wt-col-inline menu-item menu-second-skeleton-text-2 menu-second__link"></div>
                                                                    
                                            
                    <div class="wt-col-inline menu-second__buttons">
                        
                                                    <a href="/rider/download/"
                               class="menu-item menu-second__button menu-second__download-button wt-button wt-button_size_s
                                 wt-button_mode_primary">
                                Download
                            </a>
                                            </div>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
   /* site header stub height is needed to avoid Cumulative Layout Shift (CLS), which is a Web Vital */
  .site-header-stub {
    height: var(--site-header-height, 72px);
    background-color: var(--site-header-bg-color, var(--wt-color-dark));
  }

  @media (max-width: 1000px) {
    .site-header-stub--adaptive {
      height: var(--mobile-site-header-height, 48px);
    }
  }
</style>

        </div>
    
    <div class="page__content " data-js-crawler="content-root">
            <div id="react-download-thanks"></div>
        </div>

            <div class="page__footer" id="footer-container">
    <footer class="footer" id="footer"></footer>
</div>
    </div>

<script>
(function () {
  function getParameterByName(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
      results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
  }

  function updateQueryStringParameter(uri, key, value) {
    var re = new RegExp("([?&])" + key + "=.*?(&|$)", "i");
    var separator = uri.indexOf('?') !== -1 ? "&" : "?";
    if (uri.match(re)) {
      return uri.replace(re, '$1' + key + "=" + value + '$2');
    }
    else {
      return uri + separator + key + "=" + value;
    }
  }

  var downloadLink = document.getElementById("download-link");
  if (downloadLink != null) {
    var platform = getParameterByName('platform');
    platform = platform != null ? platform : "windows";
    var href = downloadLink.getAttribute("href");
    var code = getParameterByName("code");

    if(code != null) {
      href = updateQueryStringParameter(href, "code", code)
    }
    href = updateQueryStringParameter(href, "platform", platform);
    downloadLink.setAttribute("href", href);
  }
})();
</script><script>
(function() {
  var STORAGE_KEY_NAME = 'firefoxDisappearedSVGWorkaround';
  var STORAGE_KEY_VALUE = '1';

  var isFirefox = /firefox/i.test(navigator.userAgent);
  if (!isFirefox || isFirefox && sessionStorage.getItem(STORAGE_KEY_NAME) === STORAGE_KEY_VALUE) {
    return;
  }

  var arrayFrom = function (arrayLike) {
    return Array.prototype.slice.call(arrayLike, 0);
  };

  function workaround() {
    var uses = document.querySelectorAll('.page svg use');
    var badNodesCount = 0;

    arrayFrom(uses).forEach(function (node) {
      var rect = node.getBoundingClientRect();
      if (rect.width === 0 && rect.height === 0)
        badNodesCount++;
    });

    if (badNodesCount === uses.length) {
      sessionStorage.setItem(STORAGE_KEY_NAME, STORAGE_KEY_VALUE);
      if (typeof dataLayer !== 'undefined')
        dataLayer.push({'firefoxDisappearedSVGWorkaround': STORAGE_KEY_VALUE});

      window.location.replace(window.location.href);
    }
  }

  window.addEventListener('DOMContentLoaded', workaround);

})();
</script>          <script src="/_assets/banner-rotator.entry.fa25ba6764daf270853a.js" type="text/javascript"></script>
      <link href="/_assets/banner-rotator.entry.1213bcca835e111c6db1.css" rel="stylesheet" type="text/css">
  </body>
  <!--comment here-->
</html><!--end-->
	
	
	`)
	println("----------------")
	println(output)
}
