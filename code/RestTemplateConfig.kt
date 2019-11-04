package xyz.gke.mattsday.drink.save.config

import org.springframework.boot.web.client.RestTemplateBuilder
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.HttpHeaders
import org.springframework.http.client.ClientHttpRequestInterceptor
import org.springframework.web.client.RestTemplate
import javax.servlet.http.HttpServletRequest

@Configuration
class RestTemplateConfig(
        private val restTemplateBuilder: RestTemplateBuilder,
        private val request: HttpServletRequest) {

    // Headers to persist through all rest calls
    companion object {
        val INCLUDE_HEADERS = arrayOf(
                "x-request-id",
                "x-b3-traceid",
                "x-b3-spanid",
                "x-b3-parentspanid",
                "x-b3-sampled",
                "x-b3-flags",
                "x-ot-span-context",
                "user-agent",
                "cookie",
                "pipeline"
        )
    }

    val traceHeaders: HttpHeaders
        get() {
            val headers = HttpHeaders()
            try {
                val e = request.headerNames
                while (e.hasMoreElements()) {
                    val header = e.nextElement()
                    if (INCLUDE_HEADERS.contains(header)) {
                        headers.add(header, request.getHeader(header))
                    }
                }
            } catch (ignored: IllegalStateException) {
                ignored.printStackTrace()
            }
            return headers
        }

    @Bean
    fun restTemplate(): RestTemplate {
        return restTemplateBuilder.additionalInterceptors(ClientHttpRequestInterceptor
                                                          { httpRequest, bytes, clientHttpRequestExecution ->
                                                              // Persist istio tracing headers
                                                              val traceHeaders = traceHeaders
                                                              httpRequest.headers.addAll(traceHeaders)
                                                              clientHttpRequestExecution.execute(httpRequest, bytes)
                                                          }).build()
    }
}