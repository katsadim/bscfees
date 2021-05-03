const siteAddress = new URL("https://bscfees.com")

module.exports = {
    siteMetadata: {
        title: "BscFees",
        titleTemplate: "%s · HODL it!",
        description: "Bsc fees online calculator",
        url: siteAddress.hostname,
        image: "/images/icon.jpg", // Path to your image you placed in the 'static' folder
        lang: 'en',
    },
    plugins: [
        {
            resolve: "gatsby-plugin-google-gtag",
            options: {
                trackingIds: ["G-ZC8THH35LZ"],
            },
        },
        {
            resolve: "gatsby-plugin-manifest",
            options: {
                name: `BscFees`,
                short_name: `BscFees`,
                start_url: `/`,
                background_color: `#f7f0eb`,
                theme_color: `#f7f0eb`,
                display: `standalone`,
                icon: "src/images/icon.png",
            },
        },
        {
            resolve: `gatsby-plugin-s3`,
            options: {
                bucketName: 'bscfeesweb',
                protocol: siteAddress.protocol.slice(0, -1),
                hostname: siteAddress.hostname,
            },
        },
        {
            resolve: `gatsby-plugin-canonical-urls`,
            options: {
                siteUrl: siteAddress.href.slice(0, -1),
            }
        },
        `gatsby-plugin-react-helmet`,
        `gatsby-plugin-offline`,
    ],
};
