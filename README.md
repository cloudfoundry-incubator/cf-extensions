# CF-Extensions

The [`cloudfoundry-incubator/cf-extensions`](https://cloudfoundry-incubator/cf-extensions) is a Golang GitHub bot that runs periodically as a CloudFoundry Golang app and updates the [`data`](/data) and [`docs`](/docs) directory of this repository. 

The `data` folder contains the JSON databases of all discovered extensions (tracked and untracked) as well as the official list of extensions that are accepted and have statuses. 

The [`docs`](/docs) folder contains generated docs for the catalog of projects. The following diagram summarizes the work done by the `cf-extensions` bot each time it runs.

![CF-Extensions bot workflow](/docs/images/cf-extensions-bot-flowchart.png?raw=true "cf-extensions bot run process")

The resulting summary list of extensions created in the [`docs`](/docs) folder is accessible from either pointing your browser to: [`cloudfoundry-incubator/cf-extensions/docs/index.md`](https://cloudfoundry-incubator/cf-extensions/docs/imdex.md) or eventually from [`ext.cloudfoundry.org`](https://ext.cloudfoundry.org).

# Getting Started: a primer

In order to be listed as a CF-Extensions project, instigators need to follow the following steps:

1. Make their repos public.
2. If the repo is not part of one of the the CFF’s supported organizations, then move it to `github.com/cloudfoundry-community`.
3. Modify their GitHub topics list to include `cf-extensions`.
4. Next time the cf-extensions bot runs it will submit an issue to your repo with instructions to create a `.cf-extensions` JSON metadata file and what values to include.

Once you are ready to submit your proposal extension to the CF-Extensions PMC then follow these additional steps:

1. Read through the CF-Extensions process [4].
2. Create a CF-Extensions proposal following the template [3].
3. Submit proposal to [`cf-dev@lists.cloudfoundry.org`](mailto:cf-dev@lists.cloudfoundry.org).
4. Take part in the monthly CF-Extensions PMC calls [5].

During all these steps, you are welcome to contact [`cf-extensions-pmc@cloudfoundry.org`](mailto:cf-extensions-pmc@cloudfoundry.org) if you have any questions. Once a proposal is submitted and accepted, it’s status will move from “Unknown” to “Proposed” and the project will automatically move from the untracked database to the tracked database and be added to the catalog of managed CF-Extensions projects.

It’s important to note that the CF-Extensions process is self-driven. That is, each proposer/instigator can chose to move projects through the process as fast (within some limits) or as slow as they chose. A project can also be retracted at any time.

# Conclusion: how to contribute and future

The CF-Extensions GitHub bot was created out of the need to reduce this logistical burden in keeping lists of GitHub projects in syncing with the community. By requiring extensions instigators to tag (add topic in GitHub) their projects and provide a little bit of metadata, we are able to dynamically construct an accurate list of all extensions to CloudFoundry and organize this list.

Some future features may include but are not limited to:

1. Providing different views of the tracked and untracked projects databases. 
2. Expanding the search of the `cf-extensions` bot to all organizations that are part of the CFF.
3. Generating emails to `cf-dev` mailing list to communicate new proposed extensions to the community along with when their statuses have changed.

If you would like to participate in maintaining and improving the CF-Extensions GitHub bot project [1], which itself is added as an extension (tool category with distributed commit), we welcome you to clone and peruse the code. All pull requests as well as issues and suggestions and bug reports are encouraged.

# References

1. https://github.com/cloudfoundry-incubator/cf-extensions
2. [CF-Extensions Projects Lists Docs](https://docs.google.com/document/d/1EqA2vdBCzEAxCrBrhYk7tNdsx0d1hFiArNTPmKvX-qs)
3. [CF-Extensions Template](https://docs.google.com/document/d/1cpyBmds7WYNLKO1qkjhCdS8bNSJjWH5MqTE-h1UCQkQ)
4. [CF-Extensions Process](https://docs.google.com/document/d/1KaYuqNbPrr23d3OsAhi0NTwBNy-LRZK-FbN3LfBgqjw)
5. [CF-Extensions Notes](https://docs.google.com/document/d/1RCMHYFQaB1oqdEKev-cVF2Rrr6qqCT9C6RmFFKmUxnI)
6. https://developer.github.com/v3/
7. https://github.com/google/go-github/
