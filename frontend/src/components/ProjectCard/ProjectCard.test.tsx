import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import ProjectCard from './ProjectCard';

describe('ProjectCard', () => {
  const defaultProps = {
    Id: 1,
    Key: 'PRJ',
    Name: 'Test Project',
    self: 'http://example.com',
    isAdded: false,
    onAction: jest.fn().mockResolvedValue(undefined),
  };

  test('отображает имя и ключ проекта', () => {
    render(<ProjectCard {...defaultProps} />);
    expect(screen.getByText('Test Project')).toBeInTheDocument();
    expect(screen.getByText('PRJ')).toBeInTheDocument();
  });

  test('кнопка добавления вызывает onAction с действием add', async () => {
    render(<ProjectCard {...defaultProps} />);
    const button = screen.getByRole('button', { name: 'Добавить' });
    fireEvent.click(button);
    expect(defaultProps.onAction).toHaveBeenCalledWith('PRJ', 'add');
  });

  test('кнопка удаления вызывает onAction с действием remove', () => {
    const props = { ...defaultProps, isAdded: true };
    render(<ProjectCard {...props} />);
    const button = screen.getByRole('button', { name: 'Удалить' });
    fireEvent.click(button);
    expect(props.onAction).toHaveBeenCalledWith('PRJ', 'remove');
  });
});

