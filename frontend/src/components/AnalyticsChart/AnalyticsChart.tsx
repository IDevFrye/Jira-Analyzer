import React, { useEffect, useState } from 'react';
import axios from 'axios';

interface AnalyticsChartProps {
  projectKey: string;
  taskId: number;
}

const AnalyticsChart: React.FC<AnalyticsChartProps> = ({ projectKey, taskId }) => {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    axios
      .get(`/api/v1/graph/get/${taskId}`, { params: { project: projectKey } })
      .then((res) => setData(res.data));
  }, [projectKey, taskId]);

  if (!data) return <div>Загрузка графика...</div>;

  return (
    <div>
      <h4>График по задаче {taskId}</h4>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
};

export default AnalyticsChart;