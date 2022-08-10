package emaiMassage

type T string

const (
	Sold_Success_Title      T = "Congratulations, your item sold!" //卖家定价出售成功
	Sold_Success_Massage    T = "You successfully sold MetaPanance#25(NFT name) for 0.22 bNEO(price) on MegaOasis.\n"
	Auction_Success_Title   T = "Congratulations, your item auctioned!" //卖家拍卖成功
	Auction_Success_Massage T = "You successfully auctioned MetaPanance#25(NFT name) for 0.22 bNEO(price) on MegaOasis."
	Bid_Success_Title       T = "Congratulations, you bought an item at auction!" //买家竞拍成功
	Bid_Success_Massage     T = "You successfully bought MetaPanance(NFT name) for 0.22 bNEO(price) at auction on MegaOasis."
	Receive_Offer_Title     T = "receive_Offer" //卖家收到offer
	Receive_Offer_Massage   T = "receive_Offer"
	Accept_Offer_Title      T = "accept_offer" //买家offer被接受
	Accept_Offer_Massage    T = "accept_offer"
	Expired_Listed_Title    T = "expired_listed" //上架过期
	Expired_Listed_Massage  T = "expired_listed"
)

func (me T) Val() string {
	return string(me)
}
