package main

import (
	"os"
)

// ToolRegistry holds the O(1) registry mapping lowercased names/aliases to detection settings.
var ToolRegistry = make(map[string]ToolConfig)

func init() {
	// Base list of all 500+ requested tools
	standardTools := []string{
		"git", "wget", "yarn", "python", "python3", "coreutils", "pkg-config", "chromedriver", "awscli", "automake",
		"youtube-dl", "libtool", "cmake", "readline", "maven", "libyaml", "autoconf", "redis", "heroku", "rbenv",
		"libksba", "pidof", "selenium-server-standalone", "carthage", "tree", "jq", "docker", "nmap", "htop", "nvm",
		"pyenv", "ansible", "sbt", "terraform", "graphviz", "zsh-completions", "mercurial", "unrar", "git-lfs",
		"the_silver_searcher", "kubernetes-cli", "bash", "dnsmasq", "mariadb", "ant", "ruby-build", "gdb", "sqlite",
		"phantomjs", "hugo", "elixir", "fish", "watch", "swiftlint", "ghostscript", "rabbitmq", "p7zip", "flow",
		"httpie", "docker-compose", "mono", "ctags", "pandoc", "bazel", "libpng", "libffi", "tig", "llvm",
		"memcached", "ack", "xz", "docker-machine", "jpeg", "sdl2", "midnight-commander", "pyenv-virtualenv",
		"capstone", "putty", "hub", "thefuck", "gnutls", "emacs", "highlight", "reattach-to-user-namespace",
		"ssh-copy-id", "jenkins", "certbot", "fzf", "dos2unix", "rust", "packer", "graphicsmagick", "libtiff",
		"cocoapods", "mtr", "doxygen", "groovy", "winetricks", "freetype", "apache-spark", "pcre", "swig", "sdl",
		"nginx", "git-flow", "pv", "gawk", "qemu", "cairo", "docker-machine-driver-xhyve", "zsh-syntax-highlighting",
		"socat", "zeromq", "cask", "autojump", "geckodriver", "tbb", "mobile-shell", "gettext", "libxslt", "glib",
		"kotlin", "gdbm", "hadoop", "eigen", "kubernetes-helm", "exiftool", "aria2", "libevent", "binutils",
		"zlib", "bison", "rename", "snappy", "poppler", "lynx", "erlang", "mitmproxy", "openvpn", "tor", "cowsay",
		"icu4c", "aws-elasticbeanstalk", "ios-webkit-debug-proxy", "jenv", "ripgrep", "cloc", "sdl_image", "gnu-sed",
		"mas", "parallel", "glog", "speedtest-cli", "shellcheck", "sdl_mixer", "bower", "sdl_ttf", "glide",
		"postgis", "ntfs-3g", "kafka", "portmidi", "mysql-connector-c", "irssi", "leveldb", "exercism", "fswatch",
		"yasm", "kibana", "libdvdcss", "gflags", "aircrack-ng", "gd", "pango", "portaudio", "httrack", "boot2docker",
		"mcrypt", "cassandra", "clang-format", "nodebrew", "direnv", "s3cmd", "perl", "libusb", "gdal", "leiningen",
		"git-flow-avh", "webp", "terminal-notifier", "fortune", "openssh", "lmdb", "grafana", "fontconfig", "sox",
		"syncthing", "open-mpi", "colordiff", "pwgen", "swagger-codegen", "md5sha1sum", "consul", "guetzli", "iperf",
		"watchman", "libimobiledevice", "glfw", "rsync", "gtk+", "libmagic", "gifsicle", "sl", "vault", "sshfs",
		"libav", "libunistring", "influxdb", "iperf3", "scipy", "netcat", "peco", "haskell-stack", "unison", "nano",
		"pstree", "logstash", "gstreamer", "glew", "libgcrypt", "sdl2_image", "rpm", "xcproj", "gsl", "opam",
		"unixodbc", "sqlmap", "thrift", "mplayer", "z", "mosquitto", "libzip", "gtk+3", "nkf", "openexr",
		"git-extras", "ninja", "dpkg", "lftp", "bash-git-prompt", "sdl2_mixer", "pass", "proxychains-ng", "mutt",
		"media-info", "openconnect", "minicom", "npth", "sdl2_ttf", "ccache", "screen", "neovim", "ncdu", "make",
		"geoip", "nasm", "giflib", "vagrant-completion", "ruby-install", "harfbuzz", "fontforge", "brew-cask-completion",
		"upx", "apktool", "rlwrap", "lrzsz", "librsvg", "szip", "wxpython", "moreutils", "ghc", "ocaml", "valgrind",
		"chruby", "libiconv", "atool", "screenfetch", "adns", "diff-so-fancy", "siege", "apr-util", "figlet",
		"zookeeper", "libuv", "crystal-lang", "grc", "docker-completion", "curl", "mysql", "sbcl", "gpg-agent",
		"fftw", "mycli", "isl", "asciinema", "w3m", "libmemcached", "aspell", "rethinkdb", "postgresql",
		"little-cms2", "xctool", "pypy", "opus", "neo4j", "libvorbis", "apr", "dockutil", "pyenv-virtualenvwrapper",
		"htop-osx", "iftop", "ideviceinstaller", "wrk", "lame", "spark", "smpeg", "ctop", "cmatrix", "gnu-tar",
		"plantuml", "sphinx-doc", "trash", "neofetch", "pngquant", "md5deep", "optipng", "cabextract", "findutils",
		"you-get", "hydra", "wxmac", "libsodium", "pgcli", "csshx", "grep", "autoenv", "antigen", "lolcat",
		"gedit", "shadowsocks-libev", "swi-prolog", "axel", "typesafe-activator", "gnu-getopt", "smartmontools",
		"n", "e2fsprogs", "iproute2mac", "dfu-util", "pygtk", "gobject-introspection", "libssh2", "casperjs",
		"ios-deploy", "handbrake", "oniguruma", "etcd", "tidy-html5", "flex", "hashcat", "uncrustify",
		"bash-completion", "sip", "libusb-compat", "swiftformat", "pinentry-mac", "ranger", "geos", "scons",
		"mackup", "shared-mime-info", "jpegoptim", "solr", "tcptraceroute", "tesseract", "gdk-pixbuf", "libplist",
		"tmate", "weechat", "sshuttle", "nodenv", "subversion", "zsh-autosuggestions", "gource", "wakeonlan",
		"sonar-scanner", "azure-cli", "bfg", "ext4fuse", "webpack", "openshift-cli", "intltool", "msgpack",
		"markdown", "luajit", "gnuradio", "cscope", "cabal-install", "unzip", "algol68g", "mecab", "binwalk",
		"cpanminus", "autossh", "activemq", "typescript", "ipcalc", "jpeg-turbo", "byobu", "haproxy", "hive",
		"osquery", "augeas", "qcachegrind", "docker-compose-completion", "jasper", "pinentry", "jmeter",
		"arp-scan", "fabric", "links", "rbenv-gemset", "jansson", "archey", "mecab-ipadic", "libsndfile",
		"ossp-uuid", "sonarqube", "ncurses", "couchdb", "git-crypt", "fdupes", "hping", "transmission",
		"diffutils", "jemalloc", "grpc", "pigz", "grails", "dbus", "ispell", "terragrunt", "zenity", "node-build",
		"x264", "avrdude", "mpg123", "flyway", "lz4", "xmlstarlet", "libdnet", "fping", "gnu-which", "tldr",
		"numpy", "libvpx", "libev", "dex2jar", "gzip", "boost", "clisp", "netpbm", "cppcheck", "tcpreplay",
		"caddy", "xpdf", "bazaar", "clamav", "stunnel", "fasd", "usbmuxd", "elinks", "infer", "cryptopp",
		"libogg", "gmp", "mpfr", "flac", "tcl-tk", "gitlab-ci-multi-runner", "pngcrush", "ddrescue", "zplug",
		"fcrackzip", "ldid", "kafkacat", "polipo", "swiftgen", "icdiff", "glm", "privoxy", "repo", "libmpc",
		"libassuan", "rclone", "webkit2png", "nghttp2", "elm", "xhyve", "gperftools", "oath-toolkit",
		"supervisor", "chisel", "sourcekitten", "libsass", "libgpg-error", "x265", "c-ares", "libxmlsec1",
		"reaver", "pssh", "astyle", "cvs", "goaccess", "libressl", "john", "scala", "librdkafka", "less",
		"conan", "rtmpdump", "megatools", "telegraf", "dialog", "testssl", "jetty", "mpich", "freetds",
		"libgit2", "lzlib", "wireshark", "gnuplot", "miniupnpc", "sslscan", "wdiff", "libass", "ipython",
		"mongodb", "gist", "py2cairo", "pypy3", "micro", "calc", "gpac", "percona-server", "ammonite-repl",
		"pidcat", "pre-commit", "procmail", "fdk-aac", "pandoc-citeproc", "makedepend", "qpdf", "aws-sdk-cpp",
		"aws-shell", "net-snmp", "kops", "dark-mode", "pcre2", "ettercap", "docker-machine-completion",
		"pidgin", "ruby", "squid", "mit-scheme", "task", "mpv", "saltstack", "elasticsearch", "ngrep",
		"ipmitool", "pdf2htmlex", "giter8", "m-cli", "shtool", "theora", "fltk", "radare2", "tcpflow",
		"source-highlight", "imagemagick", "gdrive", "m4", "rbenv-default-gems", "texinfo", "pixman",
		"percona-toolkit", "offlineimap", "cntlm", "asciidoc", "hyper", "libcouchbase", "macvim",
		"git-credential-manager", "re2c", "lnav", "tmux", "lua", "guile", "zbar", "cmus", "berkeley-db",
		"git-standup", "jruby", "arping", "imagesnap", "grunt-cli", "libvirt", "stow", "mysql-utilities",
		"godep", "hbase", "zopfli", "expect", "hunspell", "filebeat", "v8", "libpqxx", "tomcat", "nikto",
		"cloog", "googler", "libmicrohttpd", "autoconf-archive", "libtasn1", "expat", "vegeta", "lastpass-cli",
		"mpd", "protobuf", "rocksdb", "global", "zsh", "qrencode", "p11-kit", "arangodb", "ccat", "wine",
		"sassc", "docker-machine-nfs", "freerdp", "lcov", "lzip", "pyqt", "testdisk", "gnupg", "libpcap",
		"multitail", "duti", "libtermkey", "git-review", "grace", "jsoncpp", "autogen", "sloccount", "bind",
		"foremost", "texi2html", "gradle", "doctl", "gst-plugins-good", "nuget", "prometheus", "gauge",
		"enca", "qt", "gst-plugins-bad", "vim", "purescript", "lzo", "ncftp", "atk", "mkvtoolnix",
		"freeglut", "git-cola", "ta-lib", "libtensorflow", "assimp", "grip", "minio", "ffmpeg", "libxml2",
		"flake8", "sysdig", "emscripten", "rdesktop", "planck", "dfu-programmer", "pygobject3", "ipfs",
		"openjpeg", "open-ocd", "platformio", "gnu-indent", "sfml", "rhino", "dmd", "minimal-racket",
		"sysbench", "git-quick-stats", "go", "gcc", "ios-sim", "git-town", "rpm2cpio", "gibo", "rbenv-bundler",
		"python", "cgal", "boost-python", "ghq", "libconfig", "ed", "openssl", "sphinx", "mvnvm",
		"amazon-ecs-cli", "node",
	}

	// Initialize defaults for all standard tools
	for _, tool := range standardTools {
		ToolRegistry[tool] = ToolConfig{
			Names:       []string{tool},
			VersionArgs: []string{"--version"},
			Author:      "Unknown Publisher",
			Description: "Developer tool or package",
			Example:     "`" + tool + " --help`",
		}
	}

	// Initialize top 100 NPM packages
	npmPackages := []string{
		"lodash", "chalk", "request", "commander", "react", "express", "debug", "async", "core-js", "tslib",
		"uuid", "axios", "react-dom", "mkdirp", "yargs", "glob", "colors", "inquirer", "webpack", "rxjs",
		"bluebird", "underscore", "vue", "classnames", "minimist", "body-parser", "semver", "cheerio", "eslint",
		"typescript", "js-yaml", "winston", "mocha", "socket.io", "ramda", "react-redux", "ejs", "mongoose",
		"jest", "chokidar", "nan", "postcss", "morgan", "immutable", "qs", "jsonwebtoken", "cors", "npm",
		"koa", "graphql", "prettier", "pg", "d3", "passport", "moment", "dotenv", "shelljs", "rimraf",
		"q", "minimatch", "extend", "ajv", "cross-spawn", "browserify", "execa", "prompts", "got", "node-fetch",
		"yup", "zod", "tailwindcss", "esbuild", "rollup", "vite", "prisma", "next", "nuxt", "svelte",
		"gulp", "grunt", "bower", "pnpm", "pm2", "nodemon", "concurrently", "lerna", "nx", "puppeteer",
		"cypress", "playwright", "eslint-plugin-react", "babel-core", "styled-components", "date-fns", "dayjs",
		"redux", "mobx", "chart.js", "three", "socket.io-client",
	}

	for _, pkg := range npmPackages {
		ToolRegistry[pkg] = ToolConfig{
			PackageManager: "npm",
			Author:         "Unknown Publisher",
			Description:    "NPM library package",
			Example:        "`npm install " + pkg + "`",
		}
	}

	// --- Overrides for tools requiring specific binary names, version args, or path configs ---

	ToolRegistry["java"] = ToolConfig{
		Names:       []string{"java"},
		VersionArgs: []string{"-version"},
		CommonPaths: []string{
			"/Library/Java/JavaVirtualMachines/*/Contents/Home/bin/java",
			os.Getenv("HOME") + "/.sdkman/candidates/java/*/bin/java",
		},
		Author:      "Oracle / OpenJDK",
		Description: "Platform-independent object-oriented runtime environment",
		Example:     "`java -jar app.jar` or `java Main.java`",
	}

	ToolRegistry["python"] = ToolConfig{
		Names:       []string{"python", "python3", "python2"},
		VersionArgs: []string{"--version"},
		CommonPaths: []string{
			"/usr/bin/python3",
			"/usr/local/bin/python3",
			"/opt/homebrew/bin/python3",
			os.Getenv("HOME") + "/.pyenv/shims/python",
		},
		Author:      "Python Software Foundation",
		Description: "Interpreted high-level programming language",
		Example:     "`python3 main.py` or `python -m venv venv`",
	}
	ToolRegistry["python3"] = ToolRegistry["python"] // Alias mapping

	ToolRegistry["go"] = ToolConfig{
		Names:       []string{"go"},
		VersionArgs: []string{"version"},
		CommonPaths: []string{
			"/usr/local/go/bin/go",
			os.Getenv("HOME") + "/go/bin/go",
		},
		Author:      "Google / Go Authors",
		Description: "Open-source programming language for scaling software",
		Example:     "`go run .` or `go test ./...`",
	}

	ToolRegistry["ruby"] = ToolConfig{
		Names:       []string{"ruby"},
		VersionArgs: []string{"--version"},
		CommonPaths: []string{
			"/usr/bin/ruby",
			"/usr/local/bin/ruby",
			"/opt/homebrew/bin/ruby",
			os.Getenv("HOME") + "/.rvm/rubies/*/bin/ruby",
			os.Getenv("HOME") + "/.rbenv/shims/ruby",
		},
		Author:      "Yukihiro Matsumoto / Ruby Devs",
		Description: "Dynamic, open-source programming language focusing on simplicity",
		Example:     "`ruby app.rb` or `bundle install`",
	}

	ToolRegistry["rust"] = ToolConfig{
		Names:       []string{"rustc"},
		VersionArgs: []string{"--version"},
		CommonPaths: []string{
			os.Getenv("HOME") + "/.cargo/bin/rustc",
			"/usr/local/bin/rustc",
		},
		Author:      "Rust Foundation",
		Description: "Language empowering software builders with speed and safety",
		Example:     "`cargo run` or `cargo build --release`",
	}

	ToolRegistry["node"] = ToolConfig{
		Names:       []string{"node", "nodejs"},
		VersionArgs: []string{"-v"},
		CommonPaths: []string{
			"/usr/local/bin/node",
			"/opt/homebrew/bin/node",
			os.Getenv("HOME") + "/.nvm/versions/node/*/bin/node",
		},
		Author:      "Node.js Foundation / OpenJS",
		Description: "Asynchronous event-driven JavaScript runtime environment",
		Example:     "`node app.js` or `npm run dev`",
	}

	ToolRegistry["awscli"] = ToolConfig{
		Names:       []string{"aws"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["kubernetes-cli"] = ToolConfig{
		Names:       []string{"kubectl"},
		VersionArgs: []string{"version", "--client"},
	}

	ToolRegistry["the_silver_searcher"] = ToolConfig{
		Names:       []string{"ag"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["ripgrep"] = ToolConfig{
		Names:       []string{"rg"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["cocoapods"] = ToolConfig{
		Names:       []string{"pod"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["apache-spark"] = ToolConfig{
		Names:       []string{"spark-submit", "spark-shell"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["gnu-sed"] = ToolConfig{
		Names:       []string{"gsed"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["mysql-connector-c"] = ToolConfig{
		Names:       []string{"mysql_config"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["git-flow-avh"] = ToolConfig{
		Names:       []string{"git-flow"},
		VersionArgs: []string{"version"},
	}

	ToolRegistry["open-mpi"] = ToolConfig{
		Names:       []string{"mpirun"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["md5sha1sum"] = ToolConfig{
		Names:       []string{"md5sum", "sha1sum"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["libimobiledevice"] = ToolConfig{
		Names:       []string{"ideviceinfo"},
		VersionArgs: []string{"-v"},
	}

	ToolRegistry["haskell-stack"] = ToolConfig{
		Names:       []string{"stack"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["crystal-lang"] = ToolConfig{
		Names:       []string{"crystal"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["mycli"] = ToolConfig{
		Names:       []string{"mycli"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["little-cms2"] = ToolConfig{
		Names:       []string{"jpgicc"},
		VersionArgs: []string{"-v"},
	}

	ToolRegistry["ideviceinstaller"] = ToolConfig{
		Names:       []string{"ideviceinstaller"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["gnu-tar"] = ToolConfig{
		Names:       []string{"gtar"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["sphinx-doc"] = ToolConfig{
		Names:       []string{"sphinx-build"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["shadowsocks-libev"] = ToolConfig{
		Names:       []string{"ss-local", "ss-server"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["gnu-getopt"] = ToolConfig{
		Names:       []string{"getopt"},
		VersionArgs: []string{"-V"},
	}

	ToolRegistry["iproute2mac"] = ToolConfig{
		Names:       []string{"ip"},
		VersionArgs: []string{"-V"},
	}

	ToolRegistry["ios-deploy"] = ToolConfig{
		Names:       []string{"ios-deploy"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["tidy-html5"] = ToolConfig{
		Names:       []string{"tidy"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["azure-cli"] = ToolConfig{
		Names:       []string{"az"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["openshift-cli"] = ToolConfig{
		Names:       []string{"oc"},
		VersionArgs: []string{"version"},
	}

	ToolRegistry["cabal-install"] = ToolConfig{
		Names:       []string{"cabal"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["cpanminus"] = ToolConfig{
		Names:       []string{"cpanm"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["typescript"] = ToolConfig{
		Names:       []string{"tsc"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["node-build"] = ToolConfig{
		Names:       []string{"node-build"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["gnu-which"] = ToolConfig{
		Names:       []string{"gwhich"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["dex2jar"] = ToolConfig{
		Names:       []string{"d2j-dex2jar"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["gitlab-ci-multi-runner"] = ToolConfig{
		Names:       []string{"gitlab-runner"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["kafkacat"] = ToolConfig{
		Names:       []string{"kafkacat", "kcat"},
		VersionArgs: []string{"-V"},
	}

	ToolRegistry["oath-toolkit"] = ToolConfig{
		Names:       []string{"oathtool"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["sourcekitten"] = ToolConfig{
		Names:       []string{"sourcekitten"},
		VersionArgs: []string{"version"},
	}

	ToolRegistry["pypy3"] = ToolConfig{
		Names:       []string{"pypy3"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["ammonite-repl"] = ToolConfig{
		Names:       []string{"amm"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["aws-shell"] = ToolConfig{
		Names:       []string{"aws-shell"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["giter8"] = ToolConfig{
		Names:       []string{"g8"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["radare2"] = ToolConfig{
		Names:       []string{"r2"},
		VersionArgs: []string{"-v"},
	}

	ToolRegistry["imagemagick"] = ToolConfig{
		Names:       []string{"magick", "convert"},
		VersionArgs: []string{"-version"},
	}

	ToolRegistry["git-credential-manager"] = ToolConfig{
		Names:       []string{"git-credential-manager"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["grunt-cli"] = ToolConfig{
		Names:       []string{"grunt"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["mysql-utilities"] = ToolConfig{
		Names:       []string{"mysqluc"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["lastpass-cli"] = ToolConfig{
		Names:       []string{"lpass"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["freerdp"] = ToolConfig{
		Names:       []string{"xfreerdp"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["gnupg"] = ToolConfig{
		Names:       []string{"gpg"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["purescript"] = ToolConfig{
		Names:       []string{"purs"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["git-cola"] = ToolConfig{
		Names:       []string{"git-cola"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["planck"] = ToolConfig{
		Names:       []string{"plk"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["platformio"] = ToolConfig{
		Names:       []string{"pio"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["gnu-indent"] = ToolConfig{
		Names:       []string{"gindent"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["git-quick-stats"] = ToolConfig{
		Names:       []string{"git-quick-stats"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["git-town"] = ToolConfig{
		Names:       []string{"git-town"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["amazon-ecs-cli"] = ToolConfig{
		Names:       []string{"ecs-cli"},
		VersionArgs: []string{"--version"},
	}

	ToolRegistry["swift"] = ToolConfig{
		Names:       []string{"swift"},
		VersionArgs: []string{"--version"},
		Author:      "Apple Inc.",
		Description: "Statically typed programming language for Apple platforms",
		Example:     "`swift build` or `swift test`",
	}

	ToolRegistry["xcodebuild"] = ToolConfig{
		Names:       []string{"xcodebuild"},
		VersionArgs: []string{"-version"},
		Author:      "Apple Inc.",
		Description: "Build Xcode projects and workspaces",
		Example:     "`xcodebuild -workspace MyApp.xcworkspace -scheme MyApp`",
	}

	ToolRegistry["xcrun"] = ToolConfig{
		Names:       []string{"xcrun"},
		VersionArgs: []string{"--version"},
		Author:      "Apple Inc.",
		Description: "Run Xcode developer tools",
		Example:     "`xcrun swift`",
	}

	ToolRegistry["simctl"] = ToolConfig{
		Names:       []string{"simctl"},
		VersionArgs: []string{"help"},
		Author:      "Apple Inc.",
		Description: "Control iOS Simulator",
		Example:     "`xcrun simctl list`",
	}

	ToolRegistry["fastlane"] = ToolConfig{
		Names:       []string{"fastlane"},
		VersionArgs: []string{"--version"},
		Author:      "Google",
		Description: "App automation for iOS/Android",
		Example:     "`fastlane beta`",
	}

	ToolRegistry["tuist"] = ToolConfig{
		Names:       []string{"tuist"},
		VersionArgs: []string{"version"},
		Author:      "Tuist",
		Description: "Xcode project generation tool",
		Example:     "`tuist generate`",
	}

	ToolRegistry["swiftpm"] = ToolConfig{
		Names:       []string{"swift-package", "swift"},
		VersionArgs: []string{"package", "--version"},
		Author:      "Apple Inc.",
		Description: "Swift Package Manager tool",
		Example:     "`swift package init`",
	}

	ToolRegistry["xcode-select"] = ToolConfig{
		Names:       []string{"xcode-select"},
		VersionArgs: []string{"--version"},
		Author:      "Apple Inc.",
		Description: "Manage active developer directory for Xcode",
		Example:     "`xcode-select -p`",
	}

	ToolRegistry["mint"] = ToolConfig{
		Names:       []string{"mint"},
		VersionArgs: []string{"version"},
		Author:      "Mint",
		Description: "Executable manager for Swift command line tools",
		Example:     "`mint run Carthage`",
	}

	ToolRegistry["sourcery"] = ToolConfig{
		Names:       []string{"sourcery"},
		VersionArgs: []string{"--version"},
		Author:      "Krzysztof Zabłocki",
		Description: "Code generator for Swift language",
		Example:     "`sourcery --sources ./Sources`",
	}

	// Emoji aliases for common tools
	ToolRegistry["🐍"] = ToolRegistry["python"]
	ToolRegistry["🦀"] = ToolRegistry["rust"]
	ToolRegistry["☕"] = ToolRegistry["java"]
	ToolRegistry["🐳"] = ToolRegistry["docker"]
	ToolRegistry["🐹"] = ToolRegistry["go"]
	ToolRegistry["🦫"] = ToolRegistry["go"]
	ToolRegistry["💎"] = ToolRegistry["ruby"]
}
