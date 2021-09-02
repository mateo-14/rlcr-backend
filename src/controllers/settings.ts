import { Request, Response } from 'express';
import settingsService from '../services/settings';

export const getSettings = (req: Request, res: Response) => {
  settingsService
    .getSettings()
    .then((settings) => res.json(settings))
    .catch(() => res.sendStatus(404));
};
