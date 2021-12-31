module.exports = {
   theme: "creata",
   title: "Augusteum Core",
   // locales: {
   //   "/": {
   //     lang: "en-US"
   //   },
   //   "/ru/": {
   //     lang: "ru"
   //   }
   // },
   base: process.env.VUEPRESS_BASE,
   themeConfig: {
      repo: "creatachain/augusteum",
      docsRepo: "creatachain/augusteum",
      docsDir: "docs",
      editLinks: true,
      label: "core",
      algolia: {
         id: "BH4D9OD16A",
         key: "59f0e2deb984aa9cdf2b3a5fd24ac501",
         index: "augusteum",
      },
      versions: [
         {
            label: "v0.32",
            key: "v0.32",
         },
         {
            label: "v0.33",
            key: "v0.33",
         },
         {
            label: "v0.34",
            key: "v0.34",
         },
         {
            label: "master",
            key: "master",
         },
      ],
      topbar: {
         banner: false,
      },
      sidebar: {
         auto: true,
         nav: [
            {
               title: "Resources",
               children: [
                  {
                     title: "Developer Sessions",
                     path: "/DEV_SESSIONS.html",
                  },
                  {
                     title: "RPC",
                     path: "https://docs.augusteum.com/master/rpc/",
                     static: true,
                  },
               ],
            },
         ],
      },
      gutter: {
         title: "Help & Support",
         editLink: true,
         forum: {
            title: "Augusteum Forum",
            text: "Join the Augusteum forum to learn more",
            url: "https://forum.creata.network/c/augusteum",
            bg: "#0B7E0B",
            logo: "augusteum",
         },
         github: {
            title: "Found an Issue?",
            text: "Help us improve this page by suggesting edits on GitHub.",
         },
      },
      footer: {
         question: {
            text:
               "Chat with Augusteum developers in <a href='https://discord.gg/vcExX9T' target='_blank'>Discord</a> or reach out on the <a href='https://forum.creata.network/c/augusteum' target='_blank'>Augusteum Forum</a> to learn more.",
         },
         logo: "/logo-bw.svg",
         textLink: {
            text: "augusteum.com",
            url: "https://augusteum.com",
         },
         services: [
            {
               service: "medium",
               url: "https://medium.com/@augusteum",
            },
            {
               service: "twitter",
               url: "https://twitter.com/augusteum_team",
            },
            {
               service: "linkedin",
               url: "https://www.linkedin.com/company/augusteum/",
            },
            {
               service: "reddit",
               url: "https://reddit.com/r/creatanetwork",
            },
            {
               service: "telegram",
               url: "https://t.me/creataproject",
            },
            {
               service: "youtube",
               url: "https://www.youtube.com/c/CreataProject",
            },
         ],
         smallprint:
            "The development of Augusteum Core is led primarily by [Interchain GmbH](https://interchain.berlin/). Funding for this development comes primarily from the Interchain Foundation, a Swiss non-profit. The Augusteum trademark is owned by Augusteum Inc, the for-profit entity that also maintains this website.",
         links: [
            {
               title: "Documentation",
               children: [
                  {
                     title: "Creata SDK",
                     url: "https://docs.creata.network",
                  },
                  {
                     title: "Creata Hub",
                     url: "https://hub.creata.network",
                  },
               ],
            },
            {
               title: "Community",
               children: [
                  {
                     title: "Augusteum blog",
                     url: "https://medium.com/@augusteum",
                  },
                  {
                     title: "Forum",
                     url: "https://forum.creata.network/c/augusteum",
                  },
               ],
            },
            {
               title: "Contributing",
               children: [
                  {
                     title: "Contributing to the docs",
                     url: "https://github.com/creatachain/augusteum",
                  },
                  {
                     title: "Source code on GitHub",
                     url: "https://github.com/creatachain/augusteum",
                  },
                  {
                     title: "Careers at Augusteum",
                     url: "https://augusteum.com/careers",
                  },
               ],
            },
         ],
      },
   },
   plugins: [
      [
         "@vuepress/google-analytics",
         {
            ga: "UA-51029217-11",
         },
      ],
   ],
};
