import jwt from 'jsonwebtoken';

export function generateToken(userID: string, expTime: number) {
  return new Promise((resolve, reject) => {
    jwt.sign({userID}, process.env.TOKEN_SECRET!, {expiresIn: expTime}, (err, token) => {
      if (err) reject(err);
      resolve(token);
    })
  })
}

export function verify(token: string) {
  return new Promise<string>((resolve, reject) => {
    jwt.verify(token, process.env.TOKEN_SECRET!, (err, payload) => {
      if (err) reject(err);
      resolve(payload?.userID)
    })
  })
}
