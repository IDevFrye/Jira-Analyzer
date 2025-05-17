export interface Project {
  Id: number;
  Key: string;
  Name: string;
  Url: string;
  Existence?: boolean;
}

export interface ProjectStats extends Project {
  allIssuesCount: number;
  openIssuesCount: number;
  closeIssuesCount: number;
  resolvedIssuesCount: number;
  reopenedIssuesCount: number;
  progressIssuesCount: number;
  averageTime: number;
  averageIssuesCount: number;
}

export interface Task {
  Id: number;
  Key: string;
  Name: string;
}
