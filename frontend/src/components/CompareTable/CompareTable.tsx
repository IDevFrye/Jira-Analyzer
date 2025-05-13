import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './CompareTable.scss';

interface CompareTableProps {
  selectedProjectKeys: string[];
  onSelect: (key: string) => void;
}

const CompareTable: React.FC<CompareTableProps> = ({ selectedProjectKeys, onSelect }) => {
  const [data, setData] = useState<any[]>([]);

  useEffect(() => {
    if (selectedProjectKeys.length >= 2 && selectedProjectKeys.length <= 3) {
      axios
        .get(`/api/v1/compare/1`, { params: { project: selectedProjectKeys.join(',') } })
        .then((res) => setData(res.data));
    }
  }, [selectedProjectKeys]);

  return (
    <div className="compare-container">
      <h3 className="compare-title">Сравнение проектов</h3>
      {selectedProjectKeys.length < 2 || selectedProjectKeys.length > 3 ? (
        <p className="compare-message">Выберите от 2 до 3 проектов для сравнения</p>
      ) : (
        <div className="compare-table-wrapper">
          <table className="compare-table">
            <thead>
              <tr>
                <th className="compare-header">Метрика</th>
                {selectedProjectKeys.map((key) => (
                  <th key={key} className="compare-header">{key}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {data.map((metric: any) => (
                <tr key={metric.name} className="compare-row">
                  <td className="compare-metric">{metric.name}</td>
                  {selectedProjectKeys.map((key) => (
                    <td key={key} className="compare-value">{metric[key]}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default CompareTable;