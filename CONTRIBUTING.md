# Contributing

## How do I... <a name="toc"></a>

* [Use This Guide](#introduction)?
* Ask or Say Something? 🤔🐛😱
  * [Request Support](#request-support)
  * [Report an Error or Bug](#report-an-error-or-bug)
  * [Request a Feature](#request-a-feature)
* Make Something? 🤓👩🏽‍💻📜🍳
  * [Project Setup](#project-setup)
  * [Contribute Documentation](#contribute-documentation)
  * [Contribute Code](#contribute-code)
* Manage Something ✅🙆🏼💃👔
  * [Provide Support on Issues](#provide-support-on-issues)
  * [Label Issues](#label-issues)
  * [Clean Up Issues and PRs](#clean-up-issues-and-prs)
  * [Review Pull Requests](#review-pull-requests)
  * [Merge Pull Requests](#merge-pull-requests)
  * [Tag a Release](#tag-a-release)
  * [Join the Project Team](#join-the-project-team)
* Add a Guide Like This One [To My Project](#attribution)? 🤖😻👻

## Introduction

Thank you so much for your interest in contributing!. All types of contributions are encouraged and valued. See the [table of contents](#toc) for different ways to help and details about how this project handles them!📝

## Request Support

If you have a question about this project, how to use it, or just need clarification about something:

* Open an Issue at https://github.com/SimonWaldherr/fsagent/issues
* Provide as much context as you can about what you're running into.
* Provide as many Logs from the application as possible (of course be sure to remove all confidential informations like passwords)
* Provide project and platform versions (OS, Golang, ... etc), depending on what seems relevant. If not, please be ready to provide that information if maintainers ask for it.
* Be patient - this is only a hobby project, I work more then 10 hours in my day job - after work I help in a volunteer fire department as deputy commander.
* If your problem is resolved, be nice and helpful and inform us about the solution. This helps others who have possibly the same or similar problems.

## Report an Error or Bug

If you run into an error or bug with the project:

* Open an Issue at https://github.com/SimonWaldherr/fsagent/issues
* Include *reproduction steps* that someone else can follow to recreate the bug or error on their own.
* Provide as many Logs from the application as possible (of course be sure to remove all confidential informations like passwords)
* Provide project and platform versions (OS, Golang, ... etc), depending on what seems relevant. If not, please be ready to provide that information if maintainers ask for it.

## Request a Feature

If the project doesn't do something you need or want it to do:

* Open an Issue at https://github.com/SimonWaldherr/fsagent/issues
* Provide as much context as you can about what you're running into.
* Please try and be clear about why existing features and alternatives would not work for you.

Note: It is unlikely to be able to accept every single feature request that could filed. Please understand if I need to say no.

## Project Setup

So you wanna contribute some code! That's great! This project uses GitHub Pull Requests to manage contributions, so [read up on how to fork a GitHub project and file a PR](https://guides.github.com/activities/forking) if you've never done it before.

If this seems like a lot or you aren't able to do all this setup, you might also be able to [edit the files directly](https://help.github.com/articles/editing-files-in-another-user-s-repository/) without having to do any of this setup. Yes, [even code](#contribute-code).

If you want to go the usual route and run the project locally, though:

* [Install Go](https://golang.org/doc/install)
* [Fork the project](https://guides.github.com/activities/forking/#fork)

Then in your terminal:
```sh
go get simonwaldherr.de/go/fsagent
```

And you should be ready to go!

## Contribute Documentation

Documentation is a super important part of this project and the whole Free and Open Source Community. Docs are how we keep track of what we're doing, how, and (most of all) why. It's how we stay on the same page. And it's how we tell others everything they need in order to be able to use this project -- or contribute to it. So thank you in advance.

Documentation contributions of any size are welcome! Feel free to file a PR even if you're just rewording a sentence to be more clear, or fixing a spelling mistake!

To contribute documentation:

* [Set up the project](#project-setup).
* Edit or add any relevant documentation.
* Make sure your changes are formatted correctly and consistently with the rest of the documentation.
* Re-read what you wrote, and run a spellchecker on it to make sure you didn't miss anything.
* In your commit message(s), begin the first line with `docs: `. For example: `docs: Adding a doc contrib section to CONTRIBUTING.md`.
* Write clear, concise commit message(s) using [conventional-changelog format](https://github.com/conventional-changelog/conventional-changelog-angular/blob/master/convention.md). Documentation commits should use `docs(<component>): <message>`.
* Go to https://github.com/SimonWaldherr/fsagent/pulls and open a new pull request with your changes.
* If your PR is connected to an open issue, add a line in your PR's description that says `Fixes: #123`, where `#123` is the number of the issue you're fixing.

Once you've filed the PR:

* One or more maintainers will use GitHub's review feature to review your PR.
* If the maintainer asks for any changes, edit your changes, push, and ask for another review.
* If the maintainer decides to pass on your PR, they will thank you for the contribution and explain why they won't be accepting the changes. That's ok! We still really appreciate you taking the time to do it, and we don't take that lightly. 💚
* If your PR gets accepted, it will be marked as such, and merged into the `latest` branch soon after. Your contribution will be distributed to the masses next time the maintainers [tag a release](#tag-a-release)

## Contribute Code

We like code commits a lot! They're super handy, and they keep the project going and doing the work it needs to do to be useful to others.

Code contributions of just about any size are welcome!

The main difference between code contributions and documentation contributions is that contributing code can and should contain [relevant tests](https://en.wikipedia.org/wiki/Continuous_testing) or at least informations on how to test, for the code being added or changed. 

To contribute code:

* [Set up the project](#project-setup).
* Make any necessary changes to the source code.
* Include any [additional documentation](#contribute-documentation) the changes might need.
* Write tests that verify that your contribution works as expected.
* Write clear, concise commit message(s) using [conventional-changelog format](https://github.com/conventional-changelog/conventional-changelog-angular/blob/master/convention.md).
* Dependency updates, additions, or removals must be in individual commits, and the message must use the format: `<prefix>(deps): PKG@VERSION`, where `<prefix>` is any of the usual `conventional-changelog` prefixes, at your discretion.
* Go to https://github.com/SimonWaldherr/fsagent/pulls and open a new pull request with your changes.
* If your PR is connected to an open issue, add a line in your PR's description that says `Fixes: #123`, where `#123` is the number of the issue you're fixing.

## Provide Support on Issues

Helping out other users with their questions is a really awesome way of contributing to any community. It's not uncommon for most of the issues on an open source projects being support-related questions by users trying to understand something they ran into, or find their way around a known bug.

Sometimes, the `support` label will be added to things that turn out to actually be other things, like bugs or feature requests. In that case, suss out the details with the person who filed the original issue, add a comment explaining what the bug is, and change the label from `support` to `bug` or `feature`. If you can't do this yourself, @mention a maintainer so they can do it.

In order to help other folks out with their questions:

* Go to the issue tracker and [filter open issues by the `support` label](https://github.com/SimonWaldherr/fsagent/issues?q=is%3Aopen+is%3Aissue+label%3Asupport).
* Read through the list until you find something that you're familiar enough with to give an answer to.
* Respond to the issue with whatever details are needed to clarify the question, or get more details about what's going on.
* Once the discussion wraps up and things are clarified, either close the issue, or ask the original issue filer (or a maintainer) to close it for you.

Some notes on picking up support issues:

* Avoid responding to issues you don't know you can answer accurately.
* As much as possible, try to refer to past issues with accepted answers. Link to them from your replies with the `#123` format.
* Be kind and patient with users -- often, folks who have run into confusing things might be upset or impatient. This is ok. Try to understand where they're coming from, and if you're too uncomfortable with the tone, feel free to stay away or withdraw from the issue. (note: if the user is outright hostile or is violating the CoC, [refer to the Code of Conduct](CODE_OF_CONDUCT.md) to resolve the conflict).

## Label Issues

One of the most important tasks in handling issues is labeling them usefully and accurately. All other tasks involving issues ultimately rely on the issue being classified in such a way that relevant parties looking to do their own tasks can find them quickly and easily.

In order to label issues, [open up the list of unlabeled issues](https://github.com/SimonWaldherr/fsagent/issues?q=is%3Aopen+is%3Aissue+no%3Alabel) and, **from newest to oldest**, read through each one and apply issue labels according to the table below. If you're unsure about what label to apply, skip the issue and try the next one: don't feel obligated to label each and every issue yourself!

Label | Apply When | Notes
--- | --- | ---
`bug` | Cases where the code (or documentation) is behaving in a way it wasn't intended to. | If something is happening that surprises the *user* but does not go against the way the code is designed, it should use the `enhancement` label.
`critical` | Added to `bug` issues if the problem described makes the code completely unusable in a common situation. |
`documentation` | Added to issues or pull requests that affect any of the documentation for the project. | Can be combined with other labels, such as `bug` or `enhancement`.
`duplicate` | Added to issues or PRs that refer to the exact same issue as another one that's been previously labeled. | Duplicate issues should be marked and closed right away, with a message referencing the issue it's a duplicate of (with `#123`)
`enhancement` | Added to [feature requests](#request-a-feature), PRs, or documentation issues that are purely additive: the code or docs currently work as expected, but a change is being requested or suggested. |
`help wanted` | Applied by [Committers](#join-the-project-team) to issues and PRs that they would like to get outside help for. Generally, this means it's lower priority for the maintainer team to itself implement, but that the community is encouraged to pick up if they so desire | Never applied on first-pass labeling.
`in-progress` | Applied by [Committers](#join-the-project-team) to PRs that are pending some work before they're ready for review. | The original PR submitter should @mention the team member that applied the label once the PR is complete.
`performance` | This issue or PR is directly related to improving performance. |
`refactor` | Added to issues or PRs that deal with cleaning up or modifying the project for the betterment of it. |
`starter` | Applied by [Committers](#join-the-project-team) to issues that they consider good introductions to the project for people who have not contributed before. These are not necessarily "easy", but rather focused around how much context is necessary in order to understand what needs to be done for this project in particular. | Existing project members are expected to stay away from these unless they increase in priority.
`support` | This issue is either asking a question about how to use the project, clarifying the reason for unexpected behavior, or possibly reporting a `bug` but does not have enough detail yet to determine whether it would count as such. | The label should be switched to `bug` if reliable reproduction steps are provided. Issues primarily with unintended configurations of a user's environment are not considered bugs, even if they cause crashes.
`tests` | This issue or PR either requests or adds primarily tests to the project. | If a PR is pending tests, that will be handled through the [PR review process](#review-pull-requests)
`wontfix` | Labelers may apply this label to issues that clearly have nothing at all to do with the project or are otherwise entirely outside of its scope/sphere of influence. [Committers](#join-the-project-team) may apply this label and close an issue or PR if they decide to pass on an otherwise relevant issue. | The issue or PR should be closed as soon as the label is applied, and a clear explanation provided of why the label was used. Contributors are free to contest the labeling, but the decision ultimately falls on committers as to whether to accept something or not.

## Review Pull Requests

While anyone can comment on a PR, add feedback, etc, PRs are only *approved* by team members with Issue Tracker or higher permissions.

PR reviews use [GitHub's own review feature](https://help.github.com/articles/about-pull-request-reviews/), which manages comments, approval, and review iteration.

Some notes:

* You may ask for minor changes ("nitpicks"), but consider whether they are really blockers to merging: try to err on the side of "approve, with comments".
* *ALL PULL REQUESTS* should be covered by a test: either by a previously-failing test, an existing test that covers the entire functionality of the submitted code, or new tests to verify any new/changed behavior. All tests must also pass and follow established conventions. Test coverage should not drop, unless the specific case is considered reasonable by maintainers.
* Please make sure you're familiar with the code or documentation being updated, unless it's a minor change (spellchecking, minor formatting, etc). You may @mention another project member who you think is better suited for the review, but still provide a non-approving review of your own.
* Be extra kind: people who submit code/doc contributions are putting themselves in a pretty vulnerable position, and have put time and care into what they've done (even if that's not obvious to you!) -- always respond with respect, be understanding, but don't feel like you need to sacrifice your standards for their sake, either. Just don't be a jerk about it?

## Attribution

This guide was generated using the [WeAllJS](https://wealljs.org) `CONTRIBUTING.md` generator.
If you need a `CONTRIBUTING.md` file for your Repo, use this document as inspiration or use the [weallcontribute](https://github.com/WeAllJS/weallcontribute)-cli-tool and adapt it to your needs.

