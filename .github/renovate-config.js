module.exports = {
  // Nimmt automatisch das aktuelle Repository, in dem die Action läuft
  repositories: [process.env.GITHUB_REPOSITORY],
  onboardingConfig: {
    extends: ['config:recommended'],
  },
};
