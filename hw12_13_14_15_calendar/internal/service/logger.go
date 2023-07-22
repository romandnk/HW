package service

func (s *Service) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

func (s *Service) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

func (s *Service) WriteLogInFile(path string, result string) error {
	return s.logger.WriteLogInFile(path, result)
}
