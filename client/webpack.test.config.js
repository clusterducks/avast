var path = require('path');
// webpack plugins
var ProvidePlugin = require('webpack/lib/ProvidePlugin');
var DefinePlugin    = require('webpack/lib/DefinePlugin');
var ENV = process.env.ENV = process.env.NODE_ENV = 'test';

module.exports = {
    resolve: {
        cache: false,
        extensions: [
            '',
            '.ts',
            '.js',
            '.json',
            '.css',
            '.html'
        ]
    },
    devtool: 'inline-source-map',
    module: {
        loaders: [
            // support for *.ts files.
            {
                test: /\.ts$/,
                loader: 'ts-loader',
                query: {
                    // remove TypeScript helpers to be injected below by DefinePlugin
                    compilerOptions': {
                        removeComments: true,
                        noEmitHelpers: true,
                    },
                    ignoreDiagnostics: [
                        2403, // 2403 -> Subsequent variable declarations
                        2300, // 2300 Duplicate identifier
                        2374, // 2374 -> Duplicate number index signature
                        2375  // 2375 -> Duplicate string index signature
                    ]
                },
                exclude: [
                    /\.e2e\.ts$/,
                    /node_modules/
                ]
            },
            // support for *.json files.
            {
                test: /\.json$/,
                loader: 'json-loader'
            },
            // support for *.css as raw text
            {
                test: /\.css$/,
                loader: 'raw-loader'
            },
            // support for *.html as raw text
            {
                test: /\.html$/,
                loader: 'raw-loader'
            }
        ],
        postLoaders: [
            // instrument only testing sources with istanbul
            {
                test: /\.(js|ts)$/,
                include: root('src'),
                loader: 'istanbul-instrumenter-loader',
                exclude: [
                    /\.e2e\.ts$/,
                    /node_modules/
                ]
            }
        ],
        noParse: [
            /zone\.js\/dist\/.+/,
            /angular2\/bundles\/.+/
        ]
    },

    stats: {
        colors: true,
        reasons: true
    },

    debug: false,

    plugins: [
        new DefinePlugin({
            // environment helpers
            'process.env': {
                ENV: JSON.stringify(ENV),
                NODE_ENV: JSON.stringify(ENV)
            },
            global: 'window',
            // typescript helpers
            '__metadata': 'Reflect.metadata',
            '__decorate': 'Reflect.decorate'
        }),
        new ProvidePlugin({
            // '__metadata': 'ts-helper/metadata',
            // '__decorate': 'ts-helper/decorate',
            '__awaiter': 'ts-helper/awaiter',
            '__extends': 'ts-helper/extends',
            '__param': 'ts-helper/param',
            'Reflect': 'es7-reflect-metadata/dist/browser'
        })
    ],

    // we need this due to problems with es6-shim
    node: {
        global: 'window',
        progress: false,
        crypto: 'empty',
        module: false,
        clearImmediate: false,
        setImmediate: false
    }
};

function root(args) {
    args = Array.prototype.slice.call(arguments, 0);
    return path.join.apply(path, [__dirname].concat(args));
}

function rootNode(args) {
    args = Array.prototype.slice.call(arguments, 0);
    return root.apply(path, ['node_modules'].concat(args));
}
