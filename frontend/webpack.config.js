const path = require('path');
const express = require('express');
const HtmlWebpackPlugin = require('html-webpack-plugin');

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
          filename: 'assets/images/[name][ext]' // Путь для сохранения изображений
        }
      }
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: './public/index.html',
      favicon: './public/favicon.png'
    })
  ],
  devServer: {
    port: 3000,
    hot: true,
    open: true,
    historyApiFallback: true,
    setupMiddlewares: (middlewares, devServer) => {
      if (!devServer) {
        throw new Error('webpack-dev-server is not defined');
      }

      const app = devServer.app;

      // Получение всех загруженных проектов
      app.get('/api/v1/projects', (req, res) => {
        res.json([
          { 
            Id: 1, 
            Key: 'ANLYZ', 
            Name: 'Jira Analytics', 
            Url: 'http://jira.local/browse/ANLYZ' 
          },
          { 
            Id: 2, 
            Key: 'MKTG', 
            Name: 'Marketing Campaigns', 
            Url: 'http://jira.local/browse/MKTG' 
          },
          { 
            Id: 3, 
            Key: 'DEVOPS', 
            Name: 'DevOps Tools', 
            Url: 'http://jira.local/browse/DEVOPS' 
          },
          { 
            Id: 4, 
            Key: 'CRM', 
            Name: 'CRM System', 
            Url: 'http://jira.local/browse/CRM' 
          },
          { 
            Id: 5, 
            Key: 'HRM', 
            Name: 'HR Management', 
            Url: 'http://jira.local/browse/HRM',
            Stats: true 
          },
          { 
            Id: 6, 
            Key: 'FIN', 
            Name: 'Finance Tracker', 
            Url: 'http://jira.local/browse/FIN',
            Stats: true 
          },
          { 
            Id: 7, 
            Key: 'QA', 
            Name: 'QA Automation', 
            Url: 'http://jira.local/browse/QA',
            Stats: true 
          },
          { 
            Id: 8, 
            Key: 'DOCS', 
            Name: 'Documentation Updates', 
            Url: 'http://jira.local/browse/DOCS',
            Stats: true 
          },
        ]);
      });

      // Получение сухой статистики проекта
      app.get('/api/v1/projects/:id(\\d+)', (req, res) => {
        const id = parseInt(req.params.id);
        const stats = {
          1: { Key: 'ANLYZ', Name: 'Jira Analytics', openIssuesCount: 25, closeIssuesCount: 100, resolvedIssuesCount: 80, progressIssuesCount: 5 },
          2: { Key: 'MKTG', Name: 'Marketing Campaigns', openIssuesCount: 10, closeIssuesCount: 20, resolvedIssuesCount: 10, progressIssuesCount: 5 },
          3: { Key: 'DEVOPS', Name: 'DevOps Tools', openIssuesCount: 5, closeIssuesCount: 200, resolvedIssuesCount: 150, progressIssuesCount: 30 },
          4: { Key: 'CRM', Name: 'CRM System', openIssuesCount: 8, closeIssuesCount: 12, resolvedIssuesCount: 11, progressIssuesCount: 2 },
          5: { Key: 'HRM', Name: 'HR Management', openIssuesCount: 15, closeIssuesCount: 45, resolvedIssuesCount: 40, progressIssuesCount: 5 },
          6: { Key: 'FIN', Name: 'Finance Tracker', openIssuesCount: 7, closeIssuesCount: 35, resolvedIssuesCount: 30, progressIssuesCount: 2 },
          7: { Key: 'QA', Name: 'QA Automation', openIssuesCount: 12, closeIssuesCount: 80, resolvedIssuesCount: 75, progressIssuesCount: 5 },
          8: { Key: 'DOCS', Name: 'Documentation Updates', openIssuesCount: 3, closeIssuesCount: 20, resolvedIssuesCount: 18, progressIssuesCount: 2 }
        };

        const project = stats[id] || { 
          Key: `PRJ${id}`, 
          Name: `Project ${id}`, 
          openIssuesCount: 0, 
          closeIssuesCount: 0, 
          resolvedIssuesCount: 0, 
          progressIssuesCount: 0 
        };

        res.json({
          Id: id,
          ...project,
          allIssuesCount: project.openIssuesCount + project.closeIssuesCount,
          reopenedIssuesCount: Math.floor(Math.random() * 10),
          averageTime: +(Math.random() * 10 + 5).toFixed(2), // от 5 до 15 часов
          averageIssuesCount: Math.floor(Math.random() * 20 + 5) // от 5 до 25 задач
        });
      });

      // Удаление проекта
      app.delete('/api/v1/projects/:id(\\d+)', (req, res) => {
        res.status(204).send();
      });

      // Получение доступных проектов из внешнего источника (Jira)
      app.get('/api/v1/connector/projects', (req, res) => {
        const allProjects = Array.from({ length: 50 }, (_, i) => ({
          Id: i + 1,
          Key: `EXT${i + 1}`,
          Name: `External Project ${String.fromCharCode(65 + (i % 26))}${i + 1}`,
          Url: `http://jira.local/browse/EXT${i + 1}`,
          Existence: Math.random() > 0.5
        }));

        const limit = parseInt(req.query.limit || 9);
        const page = parseInt(req.query.page || 1);
        const search = (req.query.search || '').toLowerCase();

        const filtered = allProjects.filter(p =>
          p.Key.toLowerCase().includes(search) || p.Name.toLowerCase().includes(search)
        );
        const pageCount = Math.ceil(filtered.length / limit);
        const start = (page - 1) * limit;
        const end = start + limit;

        res.json({
          Projects: filtered.slice(start, end),
          PageInfo: {
            currentPage: page,
            pageCount,
            projectsCount: filtered.length
          }
        });
      });

      // Обновление / скачивание проекта по ключу
      app.post('/api/v1/connector/updateProject', express.json(), (req, res) => {
        res.json({
          success: true,
          received: req.body,
          updatedAt: new Date().toISOString()
        });
      });

      // Получение данных по аналитической задаче
      app.get('/api/v1/graph/get/:taskNumber(\\d+)', (req, res) => {
        const taskNumber = parseInt(req.params.taskNumber);
        const data = Array.from({ length: 7 }, (_, i) => Math.floor(Math.random() * 50));
        res.json({
          taskNumber,
          project: req.query.project,
          result: data,
          labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
        });
      });

      // Проведение аналитической задачи
      app.post('/api/v1/graph/make/:taskNumber(\\d+)', express.json(), (req, res) => {
        res.json({
          taskNumber: parseInt(req.params.taskNumber),
          project: req.body.project,
          status: 'done',
          resultPreview: [1, 2, 3]
        });
      });

      // Удаление аналитических задач
      app.delete('/api/v1/graph/delete', (req, res) => {
        res.status(204).send();
      });

      // Проверка, проведена ли хотя бы одна аналитическая задача
      app.get('/api/v1/isAnalyzed', (req, res) => {
        const analyzed = Math.random() > 0.3;
        res.json({ analyzed });
      });

      // Сравнение аналитических задач между проектами
      app.get('/api/v1/compare/:taskNumber(\\d+)', (req, res) => {
        const taskNumber = parseInt(req.params.taskNumber);
        const projects = (req.query.project || '').split(',');
        const comparison = projects.map(key => ({
          project: key,
          value: Math.floor(Math.random() * 100)
        }));
        res.json({ taskNumber, comparison });
      });

      // Аналитика времени в открытом состоянии
      app.get('/api/v1/analytics/time-open', (req, res) => {
        const ranges = [
          '0-1', '1-2', '2-3', '3-5', '5-7', 
          '7-10', '10-14', '14-21', '21-30', '30+'
        ];
        
        res.json({
          project: req.query.project,
          data: ranges.map(range => ({
            range: `${range} дней`,
            count: Math.floor(Math.random() * 50) + 5
          }))
        });
      });

      // Распределение по статусам
      app.get('/api/v1/analytics/status-distribution', (req, res) => {
        const statuses = [
          { status: 'Open', color: '#ef4444' },
          { status: 'In Progress', color: '#f59e0b' },
          { status: 'Resolved', color: '#10b981' },
          { status: 'Closed', color: '#3b82f6' },
          { status: 'Reopened', color: '#8b5cf6' }
        ];
        
        res.json({
          project: req.query.project,
          data: statuses.map(s => ({
            status: s.status,
            count: Math.floor(Math.random() * 100) + 10
          }))
        });
      });

      // Затраченное время
      app.get('/api/v1/analytics/time-spent', (req, res) => {
        const users = [
          'John Doe', 'Jane Smith', 'Mike Johnson', 
          'Sarah Williams', 'David Brown', 'Emily Davis'
        ];
        
        res.json({
          project: req.query.project,
          data: users.map(user => ({
            user,
            time: Math.floor(Math.random() * 80) + 5
          }))
        });
      });

      // Распределение по приоритетам
      app.get('/api/v1/analytics/priority', (req, res) => {
        const priorities = [
          { priority: 'Critical', color: '#ef4444' },
          { priority: 'High', color: '#f97316' },
          { priority: 'Medium', color: '#f59e0b' },
          { priority: 'Low', color: '#84cc16' }
        ];
        
        res.json({
          project: req.query.project,
          data: priorities.map(p => ({
            priority: p.priority,
            count: Math.floor(Math.random() * 50) + 5
          }))
        });
      });

      // Сравнение времени в открытом состоянии
      app.get('/api/v1/compare/time-open', (req, res) => {
        const projectKeys = req.query.projects?.split(',') || [];
        
        const ranges = ['0-1', '1-2', '2-3', '3-5', '5-7', '7-10', '10-14', '14-21', '21-30', '30+'];
        
        res.json({
          projects: projectKeys,
          data: projectKeys.map(key => ({
            project: key,
            data: ranges.map(range => ({
              range: `${range} дней`,
              count: Math.floor(Math.random() * 50) + 5
            }))
          }))
        });
      });

      // Сравнение распределения по статусам
      app.get('/api/v1/compare/status-distribution', (req, res) => {
        const projectKeys = req.query.projects?.split(',') || [];
        const statuses = ['Open', 'In Progress', 'Resolved', 'Closed', 'Reopened'];
        
        res.json({
          projects: projectKeys,
          data: projectKeys.map(key => ({
            project: key,
            data: statuses.map(status => ({
              status,
              count: Math.floor(Math.random() * 100) + 10
            }))
          }))
        });
      });

      // Сравнение затраченного времени
      app.get('/api/v1/compare/time-spent', (req, res) => {
        const projectKeys = req.query.projects?.split(',') || [];
        const users = ['John Doe', 'Jane Smith', 'Mike Johnson', 'Sarah Williams', 'David Brown'];
        
        res.json({
          projects: projectKeys,
          data: projectKeys.map(key => ({
            project: key,
            data: users.map(user => ({
              user,
              time: Math.floor(Math.random() * 80) + 5
            }))
          }))
        });
      });

      // Сравнение по приоритетам
      app.get('/api/v1/compare/priority', (req, res) => {
        const projectKeys = req.query.projects?.split(',') || [];
        const priorities = ['Critical', 'High', 'Medium', 'Low'];
        
        res.json({
          projects: projectKeys,
          data: projectKeys.map(key => ({
            project: key,
            data: priorities.map(priority => ({
              priority,
              count: Math.floor(Math.random() * 50) + 5
            }))
          }))
        });
      });

      return middlewares;
    }
  }
};