import { Request, Response, NextFunction } from 'express';
import { getData } from '../services/users';

export default async function isAdmin(req: Request, res: Response, next: NextFunction) {
  const userId = req.userID;
  try {
    const userData = await getData(userId);
    if (userData?.isAdmin) {
      next();
    } else {
      res.sendStatus(401);
    }
  } catch {
    res.sendStatus(404);
  }
}
