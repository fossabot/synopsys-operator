package hub

import (
	"strings"

	"github.com/sirupsen/logrus"
)

const envOptions = `
ALLOW_ACCESS_DIRECTIVES=
BLACKDUCK_CORS_ALLOWED_HEADERS_PROP_NAME=
BLACKDUCK_CORS_ALLOWED_ORIGINS_PROP_NAME=
BLACKDUCK_CORS_EXPOSED_HEADERS_PROP_NAME=
BLACKDUCK_HUB_CORS_ENABLED=
BLACKDUCK_REPORT_IGNORED_COMPONENTS=false
BLACKDUCK_SWAGGER_DISPLAYALL=
BLACKDUCK_SWAGGER_PROXY_PREFIX=
BROKER_URL=amqps://rabbitmq/protecodesc
BROKER_USE_SSL=yes
CFSSL=cfssl:8888
CLIENT_CERT_CN=binaryscanner
DENY_ACCESS_DIRECTIVES=
DISABLE_HUB_DASHBOARD=#hub-webserver.env
HTTPS_VERIFY_CERTS=yes
HUB_LOGSTASH_HOST=logstash
HUB_POSTGRES_ADMIN=blackduck
HUB_POSTGRES_ENABLE_SSL="false"
HUB_POSTGRES_HOST=
HUB_POSTGRES_PORT=
HUB_POSTGRES_USER=blackduck_user
HUB_PROXY_DOMAIN=
HUB_PROXY_HOST=
HUB_PROXY_NON_PROXY_HOSTS=solr
HUB_PROXY_PORT=
HUB_PROXY_SCHEME=
HUB_PROXY_USER=
HUB_PROXY_WORKSTATION=
HUB_VERSION=5.0.2
HUB_WEBSERVER_PORT=8443
IPV4_ONLY=0
PUBLIC_HUB_WEBSERVER_HOST=localhost
PUBLIC_HUB_WEBSERVER_PORT=443
RABBITMQ_DEFAULT_VHOST=protecodesc
RABBITMQ_SSL_FAIL_IF_NO_PEER_CERT=false
RABBIT_MQ_HOST=rabbitmq
RABBIT_MQ_PORT=5671
SCANNER_CONCURRENCY=1
USE_ALERT=0
USE_BINARY_UPLOADS=0
~/go/src/github.com/blackducksoftware/perceptor-protoform/hack
~/go/src/github.com/blackducksoftware/perceptor-protoform/hack/hub/docker-compose ~/go/src/github.com/blackducksoftware/perceptor-protoform/hack
image: blackducksoftware/appcheck-worker:1.0.1
image: blackducksoftware/blackduck-upload-cache:1.0.2
image: blackducksoftware/hub-authentication:5.0.2
image: blackducksoftware/hub-cfssl:5.0.2
image: blackducksoftware/hub-documentation:5.0.2
image: blackducksoftware/hub-jobrunner:5.0.2
image: blackducksoftware/hub-logstash:5.0.2
image: blackducksoftware/hub-nginx:5.0.2
image: blackducksoftware/hub-postgres:5.0.2
image: blackducksoftware/hub-registration:5.0.2
image: blackducksoftware/hub-scan:5.0.2
image: blackducksoftware/hub-solr:5.0.2
image: blackducksoftware/hub-webapp:5.0.2
image: blackducksoftware/hub-zookeeper:5.0.2
image: blackducksoftware/rabbitmq:1.0.0`

// GetHubKnobs ...
func GetHubKnobs() (env map[string]string, images []string) {
	env = map[string]string{}
	images = []string{}
	logrus.Infof("%v", len(strings.Split(envOptions, "\n")))

	for _, val := range strings.Split(envOptions, "\n") {
		if strings.Contains(val, "=") {
			keyval := strings.Split(val, "=")
			env[keyval[0]] = keyval[1]
		} else if strings.Contains(val, "image") {
			fullImage := strings.Split(val, ": ")
			images = append(images, fullImage[1])
		} else {
			logrus.Infof("Skipping line %v", val)
		}
	}
	logrus.Infof("%v \n %v", images, env)
	return env, images
}
