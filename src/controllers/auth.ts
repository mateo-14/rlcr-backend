import { Request, Response } from 'express';
import * as discord from '../ds';
import { addOrUpdate, getData } from '../services/users';
import { generateToken } from '../services/tokens';

export const login = async (req: Request, res: Response) => {
  if (req.body.code) {
    try {
      const accessToken = await discord.oauth2ByCode(req.body.code);
      const user = await discord.getUserByToken(accessToken);
      await discord.addUserToGuild(user.id, accessToken);

      await addOrUpdate(user);
      const userData = await getData(user.id);

      const expireTime = parseInt(process.env.TOKEN_EXP_TIME!) * 60000;
      const token = await generateToken(user.id, expireTime);
      res.cookie('token', token, { httpOnly: true, expires: new Date(Date.now() + expireTime) });
      res.json({
        avatar: user.avatar,
        username: `${user.username}#${user.discriminator}`,
        id: user.id,
        isAdmin: userData?.isAdmin,
      });
    } catch (err) {
      console.error(err);
      res.sendStatus(401);
    }
  }
};

export const logout = (_: Request, res: Response) => {
  res.cookie('token', '', { httpOnly: true, expires: new Date(0) });
  res.sendStatus(200);
};
function oauth2ByCode() {
  throw new Error('Function not implemented.');
}
