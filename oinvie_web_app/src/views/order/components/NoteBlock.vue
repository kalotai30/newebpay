<template>
<div class="cart-list">
    <div class="li dark mb-2">
        小計
        <span class="float-right">${{originTotalPrice}}</span>
    </div>

    <div class="li" v-if="isCordova || isOnlineOrder">
        <el-button type="text" class="btn-question" @click="fareModel=true">外送費<i class="el-icon-question i-info"/></el-button>
        <span v-if="radioTake==='2'" class="float-right red">
            <span v-if="deliveryFee > 0">${{deliveryFee}}</span>
            <span v-else>免外送費</span>
        </span>
    </div>
    <!--todo 外送費說明彈框-->
    <el-dialog
            v-if="fareModel"
            :visible.sync="fareModel"
            top="20%"
            width="80%"
            append-to-body
            custom-class="commDialog">
        <div class="container">
            <h4 class="text-center">外送費說明</h4>
            <el-divider class="m-0"></el-divider>
            <div v-for="(item,index) in deliveryConditionList" :key="'deliveryConditionList'+index" class="d-flex mt-1 pb-1 border-bottom fixed-w-220">
                <div class="fixed-w-90">{{item.distance}}公里以內</div>
                <div class="flex-grow-1"><template v-if="item.deliveryFee">
                    外送費<b class="red">{{item.deliveryFee}}</b>元<br></template><template v-if="item.requiredAmount">
                    滿<b>{{item.requiredAmount}}</b>元, </template><span v-if="!item.deliveryFee" class="red">免運</span>
                </div>
            </div>

            <div v-if="maxDistance" class="blue">超過{{maxDistance}}公里無法使用店家外送</div>
            <div slot="footer" class="dialog-footer">
                <el-button type="info" @click="fareModel = false">確定</el-button>
            </div>
        </div>
    </el-dialog>

    <div class="li dark mt-2" v-if="couponChoiceListItem.length>0">優惠合計</div>
    <div v-for="(item2,index2) in couponChoiceListItem" :key="'item2'+index2" class="li">
        {{item2.name+'x'+item2.choiceAmount}}
        <span class="float-right red">
            {{
                ['','','-$'+item2.discountAmount*item2.choiceAmount,'',10 - (Number(item2.discountAmount) / 10)+'折'][item2.couponTypeId]
            }}
        </span>
    </div>

    <div v-if="isCanCoupon && (isCordova || isOnlineOrder)" class="li yellow mt-2" @click="goCoupon">
        有優惠嗎？<i class="el-icon-warning i-info"/>
    </div>
    <el-divider class="mt-3 mb-3"/>

    <!--todo 統一編號-->
    <tax-id-number v-if="isCanTaxIdNumber"
                    @update="updateTaxIdNumber"
                    :value="taxIdNumber"/>

    <h4 class="mt-3" id="PayType">
        <el-button type="text" class="btn-question" @click="payModel = true">付款方式<i class="el-icon-question i-info"/></el-button>
        <el-button type="text" @click="showBank = true" v-if="paymentType !== ''"><i class="el-icon-edit"/> 更改付款方式</el-button>
    </h4>
    <!--todo 付款方式說明彈框-->
    <el-dialog
            v-if="payModel"
            :visible.sync="payModel"
            top="20%"
            width="80%"
            :show-close="false"
            append-to-body
            custom-class="commDialog">
        <div class="container">
            <h4 class="outer-content">付款方式說明</h4>
            <el-divider class="m-0"></el-divider>
            <div class="main-content">
                {{checkoutInstructions}}
            </div>
            <div slot="footer" class="dialog-footer">
                <el-button type="info" @click="payModel = false">確定</el-button>
            </div>
        </div>
    </el-dialog>

    <div class="li" @click="showBank = true" :class="required(paymentType)">
        <span class="red" v-if="paymentType === ''"><i class="el-icon-warning"/> 請選擇付款方式</span>
        <span v-if="paymentType === '1'"><img src="@/assets/images/other_pay_cash.svg" class="card-pic">現金</span>
        <span v-if="paymentType === '2'"><img src="@/assets/images/other_pay_credit_card2.svg" class="card-pic">信用卡</span>
        <span v-if="paymentType === '3'"><img src="@/assets/images/other_pay_linepay2.svg" class="card-pic">LINE PAY</span>
        <span v-if="paymentType === '4'"><img src="@/assets/images/other_pay_easywallet2.svg" class="card-pic">悠遊付</span>
        <span v-if="paymentType === '5'"><img src="@/assets/images/other_pay_pi2.svg" class="card-pic">Pi錢包</span>
        <span v-if="paymentType === '6'"><img src="@/assets/images/other_pay_credit_card2.svg" class="card-pic">支付寶</span>
        <span v-if="paymentType === '7'"><img src="@/assets/images/other_pay_wechatpay2.svg" class="card-pic">微信支付</span>
        <span v-if="paymentType === '8'"><img src="@/assets/images/other_pay_oinpay.svg" class="card-pic">信用卡</span>
        <span class="float-right red">${{finalPrice}}</span>
    </div>

    <!--todo 高鉅信用卡彈框
    <el-dialog
            :visible.sync="isCreditCard"
            :close-on-click-modal="false"
            :show-close="false"
            custom-class="cash-dialog"
            :modal="true">
        <div class="container"> 
    -->
            <!-- 這邊放高鉅iframe 
            <iframe id="cashflow-frame" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" />
            <div slot="footer" class="dialog-footer">
                <el-button @click="isCreditCard = false" type="info">取消</el-button>
                <el-button type="primary">送出</el-button>
            </div>
        </div>
    </el-dialog>
            -->
    <template v-if="orderGift.length>0">
        <el-divider class="mt-3 mb-3"/>
        <div class="li dark">贈送優惠</div>
    </template>
    <div v-for="(item3,index3) in orderGift" :key="'item3'+index3" class="li">
        {{item3.couponListName}}
        <span class="float-right red">+{{item3.amount}} 張</span>
    </div>
</div>
</template>

<script>
    import TaxIdNumber from "../../../components/TaxIdNumber";
    import {notNull} from '@/util/common';
    export default {
        name: "NoteBlock",
        props: {
            fromRoute:{},
        },
        components:{
            TaxIdNumber,//統一編號
        },
        data() {
            return {
                originTotalPrice: this.$store.state.cartPageSnapshot.originTotalPrice,    //小記
                radioTake: this.$store.state.cartPageSnapshot.radioTake,
                fareModel: this.$store.state.cartPageSnapshot.fareModel,       //外送費說明彈框
                couponChoiceListItem: this.$store.state.cartPageSnapshot.couponChoiceListItem,   //已選擇的優惠券
                couponChoiceList: this.$store.state.cartPageSnapshot.couponChoiceList,       //已選擇的優惠券選項
                isCanCoupon:this.$store.state.cartPageSnapshot.isCanCoupon,//顯示優惠券使用欄位(開啟時，可提供顧客使用優惠券)
                isCanTaxIdNumber:this.$store.state.cartPageSnapshot.isCanTaxIdNumber,//顯示統一編號輸入欄位(開啟時，可提供顧客輸入統編)
                payModel: this.$store.state.cartPageSnapshot.payModel,        //付款說明彈框
                paymentType: this.$store.state.cartPageSnapshot.paymentType,
                isCreditCard: this.$store.state.cartPageSnapshot.isCreditCard,//高鉅信用卡彈框
                openUrl: this.$store.state.cartPageSnapshot.openUrl, //高鉅信用卡彈框網址
                orderGift: this.$store.state.cartPageSnapshot.orderGift,      //可獲得票券
                taxIdNumber: this.$store.state.cartPageSnapshot.taxIdNumber,    //統一編號
                maxDistance:this.$store.state.cartPageSnapshot.maxDistance,//最大運送距離
                couponList: this.$store.state.cartPageSnapshot.couponList,         //優惠表
                calculatorCouponPrice: this.$store.state.cartPageSnapshot.calculatorCouponPrice,    //計算經過coupon的金額
                discountLimit: this.$store.state.cartPageSnapshot.discountLimit,        //優惠限制設定
                storeId: Number(this.$route.query.storeId?this.$route.query.storeId:this.$store.state.cartStoreId),//店家id(query過來的是修改訂單)
                storeInfo:this.$store.state.cartPageSnapshot.storeInfo,
                checkoutInstructions:this.$store.state.cartPageSnapshot.checkoutInstructions,//結帳說明
                showBank: this.$store.state.cartPageSnapshot.showBank,        //更改付款方式
                finalPrice:this.$store.state.cartPageSnapshot.finalPrice,
                deliveryFee:this.$store.state.cartPageSnapshot.deliveryFee,
                deliveryConditionList:this.$store.state.cartPageSnapshot.deliveryConditionList,//外送費資訊
                concatCouponList: Number(this.$route.query.storeId) > 0,
                removeSubscribe:{},
                isCordova: this.$public.isCordova(),
                isOnlineOrder: Boolean(sessionStorage.getItem("isOnlineOrder")),
            }
        },
        computed:{
            //必填共用判斷
            required(){
                return (type)=>{
                    //不能為空值
                    return type==="" || type===undefined ? 'required':''
                }
            },
        },
        created() {
        },
        destroyed() {
          this.removeSubscribe();
        },
        mounted() {
            this.removeSubscribe = this.$store.subscribe((mutation) => {
                if (mutation.type==="updateCartOnlineOrderAndOtherSetting") {
                    this.orderDeliveryList = mutation.payload.orderDeliveryList;
                    this.userAddressData = mutation.payload.userAddressData;
                    this.mealPreparationTimeType = mutation.payload.mealPreparationTimeType;
                    this.deliveryConditionList = mutation.payload.deliveryConditionList;
                    this.storeInfo = mutation.payload.storeInfo;
                    this.maxDistance = mutation.payload.maxDistance;
                    this.isCanTaxIdNumber = mutation.payload.isCanTaxIdNumber;
                    this.checkoutInstructions = mutation.payload.checkoutInstructions;
                    let self = this;
                    //計算經過coupon計算的金額
                    this.calculatorCouponPrice = function (price, couponChoiceListItem) {
                        couponChoiceListItem.forEach((item)=>{
                            if (item.couponTypeId === 2) {
                                for (let index = 0; index < item.choiceAmount; index++) {
                                    price -= item.discountAmount;
                                    if(price < 0) price = 0
                                }
                            }
                        });

                        couponChoiceListItem.forEach((item)=>{
                        if (item.couponTypeId === 4) {
                            for (let index = 0; index < item.choiceAmount; index++) {
                                    price = self.computeType(price * ((100-item.discountAmount)/100) );
                            }
                        }
                        });
                        return price;
                    };
                    this.$store.commit("updateCartCalculatorCouponPrice",this.calculatorCouponPrice);
                }
                if (mutation.type==="updateCartFinalPrice") {
                    this.finalPrice = mutation.payload;
                }
                if (mutation.type==="updateCartOriginTotalPrice") {
                    this.originTotalPrice = mutation.payload;
                }
                if (mutation.type==="updateCartRadioTake") {
                    this.radioTake = mutation.payload;
                }
                if (mutation.type==="updateCartOnlineOrderAndOtherSetting") {
                    this.isCanCoupon = mutation.payload.isCanCoupon;
                }
                if (mutation.type==="updateCartDeliveryFee") {
                    this.deliveryFee = mutation.payload;
                }
                if (mutation.type==="updateCartPaymentType") {
                    this.paymentType = mutation.payload;
                }
                if (mutation.type==="updateCartShowBank") {
                    this.showBank = mutation.payload;
                }
                if (mutation.type==="updateCartOrderGift") {
                    this.orderGift = mutation.payload;
                }
                if (mutation.type==="updateCartCouponChoiceList") {
                    this.couponChoiceList = mutation.payload;
                }
                if (mutation.type==="updateCartcouponChoiceListItem") {
                    this.couponChoiceListItem = mutation.payload;
                }
                if (mutation.type==="updateCartTaxIdNumber") {
                  this.taxIdNumber = mutation.payload;
                }
                if (mutation.type==="updateCartDiscountLimit") {
                  this.discountLimit = mutation.payload;
                }
                if (mutation.type==="openCreditCardWindows") {
                  this.isCreditCard = mutation.payload.isCreditCard;
                  this.openUrl = mutation.payload.openUrl;
                    setTimeout(() => {
                        let frame = document.getElementById('cashflow-frame');
                        frame.src = this.openUrl;
                    }, 100);
                }
            });
        },
        updated(){
        },
        methods: {
            // 子元件統一編號回傳值
            updateTaxIdNumber(val){
                this.$store.commit("updateCartTaxIdNumber",val);
                this.taxIdNumber = val;
            },
            setupData(){
                if (notNull(this.fromRoute) && ['OrderCouponList','OrderDetail','OrderAddressList','EditOrderList',].indexOf(this.fromRoute.name) === -1) {  // 取得coupon(3)
                    this.$http.fetchWithAuth`GetUserCouponList${{
                        storeId: this.storeId,
                        guide: ['1','2','3','4','5','6'],
                    }}
                    ${json => {
                        if(this.concatCouponList === true) {
                          let tempCouponList = json.discount.concat(json.offset).filter((item)=>{return item.isUseStore;});
                          this.couponChoiceListItem.forEach((choiceEle,choiceIndex) => {
                              let forCount = 0;
                              tempCouponList.forEach((couponEle,couponIndex) => {
                                  if (choiceEle.couponListId === couponEle.couponListId && choiceEle.endTime === couponEle.endTime && choiceEle.startTime === couponEle.startTime) {
                                      tempCouponList[couponIndex].count = Number(tempCouponList[couponIndex].count) + Number(choiceEle.choiceAmount);
                                      tempCouponList[couponIndex].choiceAmount = Number(choiceEle.choiceAmount);
                                      choiceEle.code = couponEle.code;
                                      forCount += 1;
                                  }
                              });
                              if (forCount === 0) {
                                  tempCouponList.push(choiceEle);
                              }
                          });
                          this.couponList = tempCouponList;
                          this.concatCouponList = false
                        } else {
                          this.couponList = json.discount.concat(json.offset).filter((item)=>{return item.isUseStore;});
                          this.couponList.forEach((item)=>{
                              this.$set(this.couponChoiceList,item.code,0);
                          });
                          this.discountLimit = json.orderDiscountLimitSetting;
                        }
                        this.$store.commit("updateCartCouponList",this.couponList);
                        this.$store.commit("updateCartDiscountLimit",this.discountLimit);
                    }}`;
                }
            },
            computeType(number) {
              let self = this
              if (self.storeInfo.discountRole === 0) {
                return Math.round(number)
              } else if (self.storeInfo.discountRole === 1) {
                return Math.ceil(number)
              } else if (self.storeInfo.discountRole === 2) {
                return Math.floor(number)
              }
            },
            goCoupon(){
                //等0.2秒才跳頁，怕有原生鍵盤還未收回就跳頁會跑版
                setTimeout(() => {
                    this.$router.push({name: 'OrderCouponList',params:{couponList:this.couponList,couponChoiceList:this.couponChoiceList,originTotalPrice:this.originTotalPrice,storeInfo:this.storeInfo,calculatorCouponPrice:this.calculatorCouponPrice,couponChoiceListItem:this.couponChoiceListItem,discountLimit:this.discountLimit,radioTake:this.radioTake}})
                }, 200);
            },
        },
        watch:{
            showBank:{
                handler:function (newVal,oldVal) {
                    this.$store.commit("updateCartShowBank",newVal);
                },
                deep:false,
                immediate:false,
            },
        },
    }
</script>
