import { Request, Response } from 'express';
import { client } from '../ds';

export const getUser = (req: Request, res: Response) => {
  client.users
    .fetch(req.userID)
    .then((user) => {
      res.json({ avatar: user.avatar, username: `${user.username}#${user.discriminator}`, id: user.id });
    })
    .catch((err) => {
      console.error(err);
      res.sendStatus(401);
    });
};
