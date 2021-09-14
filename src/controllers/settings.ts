import { Request, Response } from 'express';
import * as settingsService from '../services/settings';

export function getSettings(req: Request, res: Response) {
  settingsService
    .getSettings()
    .then((settings) => res.json(settings))
    .catch(() => res.sendStatus(404));
}
