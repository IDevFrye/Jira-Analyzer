const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');
const Dotenv = require('dotenv-webpack');

module.exports = {
  entry: './src/index.tsx',
  output: {
    path: path.join(__dirname, '/dist'),
    filename: 'bundle.js',
    publicPath: '/'
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.scss']
  },
  module: {
    rules: [
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'babel-loader',
            options: {
              presets: [
                '@babel/preset-env',
                '@babel/preset-react',
                '@babel/preset-typescript'
              ]
            }
          },
          {
            loader: 'ts-loader',
            options: {
              compilerOptions: {
                noEmit: false
              }
            }
          }
        ]
      },
      {
        test: /\.s[ac]ss$/i,
        use: [
          'style-loader',
          'css-loader',
          'sass-loader',
        ],
      },
      {
        test: /\.(png|jpe?g|gif|svg)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'assets/images/[name][ext]'
        }
      }
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: './public/index.html',
      favicon: './public/favicon.png'
    }),
    new Dotenv(),
    new webpack.DefinePlugin({
      //Для продакшена
      'process.env.REACT_APP_API_BASE_URL': JSON.stringify(process.env.REACT_APP_API_BASE_URL || '/api')

      //'process.env.REACT_APP_API_BASE_URL': JSON.stringify(process.env.REACT_APP_API_BASE_URL || 'http://localhost:8000')
    })
  ],
  devServer: {
    //port: 3000,

    //Для продакшена
    port: 80,
    hot: true,
    open: true,
    historyApiFallback: true,
    proxy: [{
      context: ['/api'],
      target: 'http://localhost:8000',
      //Для продакшена
      //target: 'http://backend:8000',
      changeOrigin: true,
      secure: false
    }]
  }
};