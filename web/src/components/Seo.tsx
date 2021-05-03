import * as React from "react";
import {FunctionComponent} from "react";
import {Helmet} from "react-helmet"
import {useLocation} from "@reach/router"
import {graphql, useStaticQuery} from "gatsby"

interface SEOProps {
    title?: string
    titleTemplate?: string,
    description?: string,
    image?: string,
    lang?: string,
}

export const SEO: FunctionComponent<SEOProps> = ({title, titleTemplate = '%s · HODL it!', description, image, lang = 'en',}: SEOProps) => {
    const {pathname} = useLocation()
    const {site} = useStaticQuery(query)

    const {
        defaultTitle,
        defaultTitleTemplate,
        defaultDescription,
        siteUrl,
        defaultImage,
        defaultLang,
    } = site.siteMetadata

    const seo = {
        title: title || defaultTitle,
        titleTemplate: titleTemplate || defaultTitleTemplate,
        description: description || defaultDescription,
        image: `${siteUrl}${image || defaultImage}`,
        url: `${siteUrl}${pathname}`,
        lang: lang || defaultLang,
    }

    return (
        <Helmet title={seo.title} titleTemplate={titleTemplate}>
            <meta name="description" content={seo.description}/>
            <meta name="image" content={seo.image}/>

            <meta name="theme_color" content="#f7f0eb"/>

            {seo.lang && <html lang={seo.lang}/>}

            {seo.url && <meta property="og:url" content={seo.url}/>}

            {seo.title && <meta property="og:title" content={seo.title}/>}

            {seo.description && (
                <meta property="og:description" content={seo.description}/>
            )}

            {seo.image && <meta property="og:image" content={seo.image}/>}

            <meta name="twitter:card" content="summary_large_image"/>

            {seo.title && <meta name="twitter:title" content={seo.title}/>}

            {seo.description && (
                <meta name="twitter:description" content={seo.description}/>
            )}

            {seo.image && <meta name="twitter:image" content={seo.image}/>}
        </Helmet>
    )
}

const query = graphql`
  query SEO {
    site {
      siteMetadata {
        defaultTitle: title
        defaultTitleTemplate: titleTemplate
        defaultDescription: description
        siteUrl: url
        defaultImage: image
        defaultLang: lang
      }
    }
  }
`
