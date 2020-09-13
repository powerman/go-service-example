module.exports = {
    extends: ['@commitlint/config-conventional'],
    rules: {
        // Start subject with upper case to have nicer ChangeLog in GitHub Releases.
        'subject-case': [
            2,
            'always',
            ['sentence-case', 'start-case', 'pascal-case', 'upper-case'],
        ],
    },
};
