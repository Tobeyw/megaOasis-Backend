package nftEvent

type T string
const (
	Sold_Success      T = "sold_success"     //卖家定价出售成功
	Auction_Success   T = "auction_auccess"  //卖家拍卖成功
	Bid_Success       T = "bid_Success"      //买家竞拍成功
	Receive_Offer     T = "receive_Offer"    //卖家收到offer
	Accept_Offer      T = "accept_offer"     //买家offer被接受
	Expired_Listed    T = "expired_listed"             //上架过期

)

func (me T) Val() string {
	return string(me)
}