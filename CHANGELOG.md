# Changelog

## [1.1.8](https://github.com/keitap/github-asana-request-review-action/compare/v1.1.7...v1.1.8) (2026-07-08)


### Bug Fixes

* paginate story and subtask search to prevent duplicate comments ([f1ec388](https://github.com/keitap/github-asana-request-review-action/commit/f1ec388f9007fe56f7f24d77e347766d21ec46f8))
* paginate story and subtask search to prevent duplicate comments ([08b5123](https://github.com/keitap/github-asana-request-review-action/commit/08b512326f80037c224690805c6d7201c7380c61))


### CI

* bump actions/setup-go from 6.4.0 to 6.5.0 ([a7e0b3d](https://github.com/keitap/github-asana-request-review-action/commit/a7e0b3d0020b92e6256d41a3f01dc2d1099b266a))
* bump actions/setup-go from 6.4.0 to 6.5.0 ([a042140](https://github.com/keitap/github-asana-request-review-action/commit/a0421406a7f15004214898bd50c4fac39c1d03f5))
* bump docker/build-push-action from 7.2.0 to 7.3.0 ([1e72d42](https://github.com/keitap/github-asana-request-review-action/commit/1e72d4293233fba322a600a9cf240c9a2261f9f4))
* bump docker/build-push-action from 7.2.0 to 7.3.0 ([6c89f7c](https://github.com/keitap/github-asana-request-review-action/commit/6c89f7c548f9dc9f599cce36055be76ef2e86c5f))
* bump docker/setup-buildx-action from 4.1.0 to 4.2.0 ([e77a0a2](https://github.com/keitap/github-asana-request-review-action/commit/e77a0a25aef8017e64adfce476d399fc0a26aa9a))
* bump docker/setup-buildx-action from 4.1.0 to 4.2.0 ([738793c](https://github.com/keitap/github-asana-request-review-action/commit/738793c47f6de7a7f074ab08053aedc861acf566))
* bump golangci/golangci-lint-action from 9.2.1 to 9.3.0 ([e665952](https://github.com/keitap/github-asana-request-review-action/commit/e6659524f919da5f0caa317f84e037a9b4dc7dc8))
* bump golangci/golangci-lint-action from 9.2.1 to 9.3.0 ([0323285](https://github.com/keitap/github-asana-request-review-action/commit/03232856423c897b8a9f735a9971963b38dd48f4))
* bump trufflesecurity/trufflehog from 3.95.6 to 3.95.8 ([00f3fbc](https://github.com/keitap/github-asana-request-review-action/commit/00f3fbc8a88eb38819f9dce48a4f0320c4c8b2bf))
* bump trufflesecurity/trufflehog from 3.95.6 to 3.95.8 ([7fb9149](https://github.com/keitap/github-asana-request-review-action/commit/7fb914986edfe587eb3d497689b7505cdb84f6d2))

## [1.1.7](https://github.com/keitap/github-asana-request-review-action/compare/v1.1.6...v1.1.7) (2026-06-23)


### Bug Fixes

* escape Asana html_text and dedupe comment emoji ([520ebe7](https://github.com/keitap/github-asana-request-review-action/commit/520ebe7c5a3c06817cdd8694fbb98b70fe271b26))
* escape dynamic values in Asana html_text and dedupe comment emoji ([d80072b](https://github.com/keitap/github-asana-request-review-action/commit/d80072bd0a3d7a6e35367443c8917cf788840a27))
* use github.Ptr instead of deprecated github.String in test ([96f48f8](https://github.com/keitap/github-asana-request-review-action/commit/96f48f8abef9e9e8a0a2b858d51063f56d0aa37f))


### Dependencies

* bump golang.org/x/oauth2 from 0.35.0 to 0.36.0 ([eacc5f6](https://github.com/keitap/github-asana-request-review-action/commit/eacc5f61024e29ac0786d50eb1cf1b970f08e874))


### CI

* add Dependabot config for gomod, github-actions and docker ([d703b50](https://github.com/keitap/github-asana-request-review-action/commit/d703b507993d74668a638e1a6755093386923758))
* add govulncheck workflow ([6752728](https://github.com/keitap/github-asana-request-review-action/commit/6752728cd7acffd3ab0a67da939c7ed72ecca846))
* add govulncheck workflow ([923e605](https://github.com/keitap/github-asana-request-review-action/commit/923e605a95b24fe06063dac83d3fea6ea28c48a7))
* add security scanning (zizmor + trufflehog) and harden workflows ([90c1484](https://github.com/keitap/github-asana-request-review-action/commit/90c1484db08a248db30b1023568647dc5179267d))
* add security workflow with zizmor and trufflehog ([c41a10f](https://github.com/keitap/github-asana-request-review-action/commit/c41a10f48be47aa964f0a68807112aa1464a775d))
* bump actions/checkout to v7 and release-please-action to v5 ([dac09d6](https://github.com/keitap/github-asana-request-review-action/commit/dac09d6a4dcb197cfe80c052d8cf0218985239db))
* bump docker actions and golangci-lint-action ([b01928f](https://github.com/keitap/github-asana-request-review-action/commit/b01928f90720eaaa491a06fa321ad95f95a2b274))
* bump Go to 1.26 ([7c33bbc](https://github.com/keitap/github-asana-request-review-action/commit/7c33bbcc781852ee6f87ac8be1d3cf3f476a297f))
* bump Go to 1.26 and update GitHub Actions ([f6d1093](https://github.com/keitap/github-asana-request-review-action/commit/f6d1093f3328c28211f5a301a46c8eb1db9a945d))
* bump golangci-lint to v2.12.2 and ignore tests for goconst ([9155da8](https://github.com/keitap/github-asana-request-review-action/commit/9155da86398926e585d93f5829ec0d431f9cd781))
* harden the go-test job after rebase ([ff0732f](https://github.com/keitap/github-asana-request-review-action/commit/ff0732f30159704c86dd276d9462ae97bea91a31))
* pin actions to commit SHA and harden workflow permissions ([46941c9](https://github.com/keitap/github-asana-request-review-action/commit/46941c98284d7419ce9fb93968a34115f189e85b))
* run go test in the test workflow ([7ad13d3](https://github.com/keitap/github-asana-request-review-action/commit/7ad13d3dc2a77735cee2fa0ea44f5f7da6e9b32d))
* run go test in the test workflow ([4a09b10](https://github.com/keitap/github-asana-request-review-action/commit/4a09b10b1d250607319fd79aae286e1cfd63a92e))


### Tests

* add E2E coverage for blockquote review body ([8a30671](https://github.com/keitap/github-asana-request-review-action/commit/8a306711ee69b9c93b661959e360d43a2dfda6d5))
* skip Asana integration tests only when explicitly requested ([0b652b5](https://github.com/keitap/github-asana-request-review-action/commit/0b652b590c30c2070a03ecc75b5c041bb751641f))
* skip Asana integration tests only when explicitly requested ([bf28535](https://github.com/keitap/github-asana-request-review-action/commit/bf28535dc06ee8ae5e8f7881268b32a6d64048a0))

## [1.1.6](https://github.com/keitap/github-asana-request-review-action/compare/v1.1.5...v1.1.6) (2026-03-17)


### Bug Fixes

* fail action on error & add release-please ([4fd0c9e](https://github.com/keitap/github-asana-request-review-action/commit/4fd0c9e311af7fe8a5dc7d68b04256c8b7ffb7db))
* fail GitHub Action when Handle() returns an error ([ddcdafc](https://github.com/keitap/github-asana-request-review-action/commit/ddcdafc4e5922682fe85ec48bec8194dec3c9d69))


### Miscellaneous

* unhide all changelog sections ([18e30e2](https://github.com/keitap/github-asana-request-review-action/commit/18e30e268ad2f2ab1563cf5bbbbb600db95715c3))
* unhide all changelog sections in release-please config ([8c9ad21](https://github.com/keitap/github-asana-request-review-action/commit/8c9ad21647647903e18c454665512d776cbf7b55))


### CI

* add release-please for automated versioning and releases ([7a980ed](https://github.com/keitap/github-asana-request-review-action/commit/7a980ed404099fba8d80c961ca1d67c8acab803c))
* use local Dockerfile for e2e testing ([7c8b5a0](https://github.com/keitap/github-asana-request-review-action/commit/7c8b5a0a109dd7e3542455de0971d2707347a3e8))
* use local Dockerfile for e2e testing ([1548f44](https://github.com/keitap/github-asana-request-review-action/commit/1548f441e56768c4af0afd168b8a6fec101d28d4))
