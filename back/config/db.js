const mongoose = require("mongoose")

//const connectDB = async () => {
//    try {
//        mongoose.set('strictQuery', false);

//        mongoose.connect(process.env.MONGO_URI, () => {
//            console.log("mongo connecté")
//        })
//    } catch (err) {
//        console.log("echec de la connection a la db " + "\n" + err)
//        process.exit()
//    }
//}

module.exports = connectDB