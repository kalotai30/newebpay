<template>
    <div class="my-cart">
        <navigation-bar Back/>
        <div class="main">
            <!--TODO 線上點餐＋候位預約點餐 才有以下區塊-->
          <MealList :fromRoute="fromRoute" :takeoutId="takeoutId" ref="MealList"></MealList>

            <!--TODO 揪團點餐 才有以下區塊-->
            <!--        <template>
                        <div class="top-cart-list">
                            <h4>
                                我的訂單
                                <el-button type="text" @click="$router.push({name: 'OrderList'})">新增餐點</el-button>
                            </h4>

                            <user style="margin: 10px 20px 0;"></user>

                            <slip-del>
                                <meal toDetail></meal>
                                <div slot="del"><img src="@/assets/images/delete_white_line.svg"> </div>
                            </slip-del>
                            <slip-del>
                                <meal toDetail></meal>
                                <div slot="del"><img src="@/assets/images/delete_white_line.svg"> </div>
                            </slip-del>

                        </div>

                        <div class="shadow-line m-0"></div>
                        <el-collapse class="top-cart-list">
                            <el-collapse-item>
                                <h4 slot="title">
                                    看看其他人都點甚麼
                                    <span class="small">張大大等，共<b class="red">10</b>人</span>
                                </h4>
                                <div class="cart-list team-buying">
                                    <user></user>
                                    <meal></meal>
                                    <meal></meal>
                                </div>
                                <div class="cart-list team-buying">
                                    <user></user>
                                    <meal></meal>
                                </div>
                            </el-collapse-item>
                        </el-collapse>
                    </template>-->

            <!--        <div class="cart-recommend">
                        <div class="con">
                            <h4>推薦</h4>
                            &lt;!&ndash;TODO 推薦商品區塊&ndash;&gt;
                            <Recommend></Recommend>
                        </div>
                    </div>-->
            <DeliveryBlock :fromRoute="fromRoute" ref="DeliveryBlock"></DeliveryBlock>
            <NoteBlock :fromRoute="fromRoute" ref="NoteBlock"></NoteBlock>
            <div class="cart-list">
                <div class="li dark" v-if="isRequestTableware && (isCordova || isOnlineOrder)">
                    <i class="el-icon-knife-fork"/>請問你需要拋棄式餐具、吸管嗎？
                </div>
                <div class="li" v-if="isRequestTableware && (isCordova || isOnlineOrder)">
                    <span v-if="!isRequiredTableware">不用提供，我是環保小尖兵</span>
                    <span v-else>需要提供</span>
                    <span class="float-right">
                        <el-switch class="lg" v-model="isRequiredTableware"/>
                    </span>
                </div>

                <el-divider class="mt-3 mb-3"/>

                <h4>
                    給店家的話
                </h4>
                <el-input
                        type="textarea"
                        :autosize="{ minRows: 2, maxRows: 4}"
                        :placeholder="remarkPrompt"
                        v-model="remark">
                </el-input>
            </div>

            <div class="cart-list last">
                請注意：當按下「送出訂單」即表示成立訂單並同意服務條款及隱私條款。
            </div>

            <!--TODO 更改付款方式-->
            <bank-card  v-if="showBank"
                        :paymentType="paymentType"
                        :paymentMethod="paymentMethod"
                        @changePaymentType="changePaymentType"/>

            <!--todo 刪除餐點彈框-->
            <el-dialog
                    :visible.sync="showDelDialog"
                    append-to-body
                    top="20%"
                    width="80%"
                    custom-class="commDialog"
                    @close="closeDelDialog">
                <div class="container">
                    <img src="../../assets/images/dialog-order.svg">
                    <h4 class="text-center">確定要刪除此餐點？</h4>
                    <div class="dialog-footer">
                        <el-button @click="delMeal" type="info">確定刪除</el-button>
                    </div>
                </div>
            </el-dialog>

            <!--todo 點餐成功彈框-->
            <el-dialog
                    v-if="completeModel"
                    :visible.sync="completeModel"
                    title="點餐成功"
                    top="20%"
                    width="80%"
                    class="dialog-complete"
                    custom-class="commDialog"
                    @close="closeSuccessDialog">
                <div class="container text-center">
                    <div>{{completeOrderData.userName}} 您好<br>以下是您的點餐資訊</div>
                    <div v-if="completeOrderData.radioTake === 1">預計完成時間</div>
                    <div v-else-if="completeOrderData.radioTake === 2">預計抵達時間</div>
                    <div class="red date">{{completeOrderData.day}} {{completeOrderData.week}}</div>
                    <div class="time">{{completeOrderData.time}}</div>
                    <div>訂單金額 
                        <span v-if="completeOrderData.radioTake === 2">(含運)</span>
                        <span v-if="!digitalPayShowStr">(尚未支付)</span>
                    </div>
                    <h2 class="red fee"><span>$</span>{{completeOrderData.totalPrice}}</h2>

                    <!--todo 選擇外送才會有以下小短腿-->
                    <template v-if="false">
                        <h4>想要聯繫外送員?</h4>
                        <a href="http://line.me/ti/p/@pkf7762b" target="_blank" class="line">
                            <img src="../../assets/images/line_@pkf7762b.jpg">
                            小短腿Line@
                        </a>
                    </template>
                </div>
                <div slot="footer" class="dialog-footer">
                    <el-button type="info" @click="closeSuccessDialog">確定</el-button>
                </div>
            </el-dialog>

            <!--todo 線上付款失敗彈框-->
            <el-dialog
                    :visible.sync="digitalPayModel"
                    top="20%"
                    width="80%"
                    custom-class="commDialog">
                <div v-if="digitalPayShowStr" class="container text-center">
                    <div>訂單處理中...<br />請勿關閉視窗</div>
                </div>
                <div v-if="!digitalPayShowStr" class="container text-center">
                    {{digitalPayErrStr}}
                </div>
            </el-dialog>

            <!--todo 取消揪團後，被分享者要出現彈框-->
            <!--<el-dialog
                    :visible.sync="showModel"
                    top="20%"
                    width="80%"
                    :show-close="false"
                    :close-on-click-modal="false"
                    custom-class="commDialog">
                <div class="container">
                    <img src="@/assets/images/dialog-order.svg">
                    <p class="text-center">主揪已取消揪團<br>但您的訂單還是被保留可繼續下訂</p>
                    <div slot="footer" class="dialog-footer">
                        <el-button @click="showModel=false" type="info">繼續</el-button>
                    </div>
                </div>
            </el-dialog>-->
            <loading
                    :active.sync=isLoading
                    :can-cancel="false"
                    :is-full-page="true">
            </loading>
        </div>
        <!--TODO 候位預約點餐出現區塊-->
        <!--        <cart-button toRecord v-if="false"
                        :message='"送出預點清單"'/>-->

        <!--TODO 揪團點餐-主揪 出現區塊-->
        <!--        <cart-button :message="'送出訂單'" toRecord isShare isInvite isOrder isCancel></cart-button>-->
        <!--TODO 揪團點餐-被分享者 出現區塊-->
        <!--        <cart-button :message="'送出訂單'" isShare isOrder isQuit v-if="false"></cart-button>-->

        <!--TODO 送出訂單出現彈框 點餐成功彈框-->
        <cart-button
                :totalPrice="finalPrice"
                :totalAmount="sumAmount"
                :submitToDo="submitToDo"
                :message='"送出訂單"'/>
    </div>
</template>

<script>

    import CartButton from '@/views/order/components/CartButton';
    import MealList from '@/views/order/components/MealList';
    import DeliveryBlock from '@/views/order/components/DeliveryBlock';
    import NoteBlock from '@/views/order/components/NoteBlock';
    import BankCard from '@/views/order/components/BankCard';
    //import User from '@/views/order/components/User'
    //import Recommend from '@/views/order/components/Recommend'
    import {notNull} from '@/util/common';
    import * as mutationTypes from "@/store/mutationTypes";
    import * as constDefine from '@/store/constDefine';

    export default {
        name: "OrderOrderCartÍ",
        components: {
            CartButton,//底下按鈕
            MealList,
            DeliveryBlock,
            NoteBlock,
            BankCard,//更改付款方式
            //User,//揪團大頭貼
            //Recommend,//推薦商品區塊
        },
        data() {
            return {
                isLoading: false,
                orderSn: "",
                radioTake: this.$store.state.cartPageSnapshot.radioTake,         //取餐方式
                showBank: this.$store.state.cartPageSnapshot.showBank,        //更改付款方式

                paymentType: this.$store.state.cartPageSnapshot.paymentType,
                showModel: false,       //取消揪團時，被分享者要出現的提示
                completeModel: false,   //點餐成功彈框
                fareModel: this.$store.state.cartPageSnapshot.fareModel,       //外送費說明彈框
                payModel: this.$store.state.cartPageSnapshot.payModel,        //付款說明彈框
                showDelDialog: false,   //刪除餐點彈框

                scrollHeight: this.$store.state.cartPageSnapshot.scrollHeight,
                scrollTop: this.$store.state.cartPageSnapshot.scrollTop,
                storeId: Number(this.$route.query.storeId?this.$route.query.storeId:this.$store.state.cartStoreId),//店家id(query過來的是修改訂單)
                shoppingCartMealChoiceList: this.$store.state.cartPageSnapshot.shoppingCartMealChoiceList, //購物車中的餐點資訊
                shoppingCartMealData: this.$store.state.cartPageSnapshot.shoppingCartMealData,       //所有餐點資料
                refreshChoiceList: this.$store.state.cartPageSnapshot.refreshChoiceList,
                delMealIndex: null,     //暫存正在被刪除的餐點
                delMealOldVal: null,    //刪除前的數量
                delMealArr: [],         //紀錄被刪除的餐點
                originTotalPrice: this.$store.state.cartPageSnapshot.originTotalPrice,    //小記
                sumAmount: 0,           //總餐點數量
                couponList: this.$store.state.cartPageSnapshot.couponList,         //優惠表
                discountLimit: this.$store.state.cartPageSnapshot.discountLimit,        //優惠限制設定
                concatCouponList: this.$store.state.cartPageSnapshot.concatCouponList,    //是否拼接優惠表
                couponChoiceList: this.$store.state.cartPageSnapshot.couponChoiceList,       //已選擇的優惠券選項
                couponChoiceListItem: this.$store.state.cartPageSnapshot.couponChoiceListItem,   //已選擇的優惠券
                calculatorCouponPrice: this.$store.state.cartPageSnapshot.calculatorCouponPrice,    //計算經過coupon的金額
                orderGift: this.$store.state.cartPageSnapshot.orderGift,      //可獲得票券
                taxIdNumber: this.$store.state.cartPageSnapshot.taxIdNumber,    //統一編號
                remark: this.$store.state.cartPageSnapshot.remark,         //訂單備註

                storeInfo:this.$store.state.cartPageSnapshot.storeInfo,
                userAddressData:this.$store.state.cartPageSnapshot.userAddressData,
                userAddressId:this.$store.state.cartPageSnapshot.userAddressId,
                deliveryConditionList:this.$store.state.cartPageSnapshot.deliveryConditionList,//外送費資訊
                maxDistance:this.$store.state.cartPageSnapshot.maxDistance,//最大運送距離
                isRequiredRemark:this.$store.state.cartPageSnapshot.isRequiredRemark,//訂單備註欄位是否必填
                isRequestTableware:this.$store.state.cartPageSnapshot.isRequestTableware,//是否顯示需要拋棄式餐具、吸管嗎？
                isRequiredTableware:this.$store.state.cartPageSnapshot.isRequiredTableware,//是否需要拋棄式餐具、吸管
                remarkPrompt:this.$store.state.cartPageSnapshot.remarkPrompt,//訂單備註欄位提示
                checkoutInstructions:this.$store.state.cartPageSnapshot.checkoutInstructions,//結帳說明
                isCanCoupon:this.$store.state.cartPageSnapshot.isCanCoupon,//顯示優惠券使用欄位(開啟時，可提供顧客使用優惠券)
                isCanTaxIdNumber:this.$store.state.cartPageSnapshot.isCanTaxIdNumber,//顯示統一編號輸入欄位(開啟時，可提供顧客輸入統編)
                mealPreparationTimeType:this.$store.state.cartPageSnapshot.mealPreparationTimeType,//餐點平均準備時間

                orderDeliveryList:this.$store.state.cartPageSnapshot.orderDeliveryList,//取餐方式（自取、外送、小短腿）
                completeOrderData: this.$store.state.cartPageSnapshot.completeOrderData,

                isCreditCard: this.$store.state.cartPageSnapshot.isCreditCard,//高鉅信用卡彈框
                takeoutId:Number(this.$store.state.cartPageSnapshot.takeoutId?this.$store.state.cartPageSnapshot.takeoutId:(this.$route.query.takeoutId?this.$route.query.takeoutId:0)),//我想修改訂單回來帶入參數
                paymentMethod:this.$store.state.cartPageSnapshot.paymentMethod,
                orderTmpId: this.$store.state.cartPageSnapshot.orderTmpId, //暫時訂單編號
                openUrl: this.$store.state.cartPageSnapshot.openUrl,
                routerGoEnable: this.$store.state.cartPageSnapshot.routerGoEnable,
                deliveryFee: this.$store.state.cartPageSnapshot.deliveryFee,
                removeSubscribe: {},
                fromRoute:null, //(不可寫進vuex)此為上一頁的路由
                timeoutID: 0,
                checkCumulativeTypeList: {},
                orderSnOCP: '',  //OCP電子錢包回傳訂單編號用

                digitalPayModel: false,     //線上金流彈窗
                digitalPayShowStr: false,   //線上金流需顯示字串

                isCordova: this.$public.isCordova(),
                isOnlineOrder: Boolean(sessionStorage.getItem("isOnlineOrder")),
                dataObject: {},
                phone: '',
                tableId: 0,
                adult: 0,
                child: 0,
                digitalPayErrStr: "",       //線上金流錯誤訊息
            }
        },
        created() {
            if(!this.isCordova) {
              try {
                this.dataObject = this.$crypto.uriComponentToObject(sessionStorage.getItem(constDefine.StorageNameType.TableOrderDataToken));
                if(this.dataObject.isChecked) {
                  this.storeId = Number(this.dataObject.storeId);
                  this.phone = this.dataObject.phone;
                  this.tableId = Number(this.dataObject.tableId);
                  this.adult = Number(this.dataObject.adult);
                  this.child = Number(this.dataObject.child);
                } else {
                  if(this.dataObject.routePath) {
                    this.$router.push({path: this.dataObject.routePath})
                  } else {
                    this.$router.push({name:'Error'});
                  }
                  return false;
                }
              } catch (e) {
                if(!this.isOnlineOrder) {
                  this.$router.push({name:'Error'});
                  return false;
                }
              }
            }

            //等0.3秒抓beforeRouteEnter
            setTimeout(() => {
                this.$refs["MealList"].setupData();
                if(this.isCordova || this.isOnlineOrder) {
                  this.$refs["DeliveryBlock"].setupData();
                  this.$refs["NoteBlock"].setupData();
                }
                this.setupData();
            }, 300);

            this.$store.commit("trigerCartDelDialog",{
              showDelDialog: false,
              delMealIndex: null,
              delMealOldVal: null,
            });
        },
        mounted() {
          this.removeSubscribe = this.$store.subscribe((mutation, state) => {
            if(mutation.type === mutationTypes.setOrderStoreInfo) {
              this.storeInfo = state.orderStoreInfo
            }

            if (mutation.type==="updateShoppingCartMealChoiceList") {
              this.shoppingCartMealChoiceList = mutation.payload;
            }
            if (mutation.type==="trigerCartDelDialog") {
              this.showDelDialog = mutation.payload.showDelDialog;
              this.delMealIndex = mutation.payload.delMealIndex;
              this.delMealOldVal = mutation.payload.delMealOldVal;
            }
            if (mutation.type==="updateCartTaxIdNumber") {
              this.taxIdNumber = mutation.payload;
            }
            if (mutation.type==="updateCartCalculatorCouponPrice") {
              this.calculatorCouponPrice = mutation.payload;
            }
            if (mutation.type==="updateCartCouponList") {
              this.couponList = mutation.payload;
            }
            if (mutation.type==="updateCartDiscountLimit") {
              this.discountLimit = mutation.payload;
            }
            if (mutation.type==="updateCartShowBank") {
              this.showBank = mutation.payload;
            }
            if (mutation.type ==="updateCartRadioTake") {
              this.radioTake = mutation.payload;
            }
            if (mutation.type ==="updateCartDeliveryFee") {
              this.deliveryFee = mutation.payload;
            }
            if (mutation.type ==="updateCartUserAddressData") {
              this.userAddressData = mutation.payload;
            }
          })

          if (localStorage.getItem('routerGoEnable') === 1) {
            localStorage.removeItem('routerGoEnable');
            this.$router.go(-1);
          }

          let checkCumulativeTypeList = {}

          this.couponChoiceListItem.forEach(item => {
            let couponListId = "couponListId_"+item.couponListId
            if (item.posUseType === 2) {
              let newItem = JSON.parse(JSON.stringify(item))
              if (typeof checkCumulativeTypeList[couponListId] === "undefined") {
                newItem.totalChoiceAmount = item.choiceAmount;
              } else {
                newItem.totalChoiceAmount = checkCumulativeTypeList[couponListId].totalChoiceAmount + item.choiceAmount
              }
              checkCumulativeTypeList[couponListId] = newItem
            }
          })

          this.checkCumulativeTypeList = checkCumulativeTypeList

        },
        destroyed() {
            this.removeSubscribe();
            if (this.timeoutID !== 0) {
              window.clearInterval(this.timeoutID);
            }
        },
        methods: {
            setupData(){
                let mealChoiceListTable = '';   //sql_ite餐點選擇狀態資料表名稱
                let mealDataTable = '';     //sql_lite參點資料的資料表名稱
                //判斷是否為修改訂單
                if (this.takeoutId!==0) {
                    mealChoiceListTable = 'order_record_meal_choiceList'; 
                    mealDataTable = 'order_record_meal_data';
                }else{
                    mealChoiceListTable = 'shopping_cart_meal_choiceList';
                    mealDataTable = 'shopping_cart_meal_data';
                }
                // 把變數還原到離開前的狀態
                // if ((notNull(this.fromRoute) && ['OrderCouponList','OrderDetail','OrderAddressList','EditOrderList'].indexOf(this.fromRoute.name) !== -1)&&notNull(this.$store.state.pageSnapshot[this.$route.name])) {
                //         for (const key in snapshotdata) {
                //             if (Object.hasOwnProperty.call(snapshotdata, key)) {
                //                 const element = snapshotdata[key];
                //                 this[key] = element;
                                
                //             }
                //         }
                // }
                if (notNull(this.fromRoute) && ['OrderCouponList','OrderDetail','OrderAddressList'].indexOf(this.fromRoute.name) !== -1) {  // 渲染完成後等一下才捲到目標位置(5)
                    this.$nextTick(function () {
                        if (this.scrollTop !== '') {
                            let finish = false;
                            let count = 0;
                            this.timeoutID = setInterval(() => {
                                // 當高度與離開前相等捲到先前高度（有時候選擇器會選不到）
                                if (document.querySelector('.main')!==null&&(Math.abs(this.scrollHeight - document.querySelector('.main').scrollHeight)<200)) {
                                document.querySelector('.main').scrollTo(0,this.scrollTop);
                                this.scrollTop = '';
                                finish = true; 
                                }
                                count = count + 1;
                                if (finish || count>200) {
                                    window.clearInterval(this.timeoutID);
                                }
                            }, 20);
                        }
                    });
                }
                if (notNull(this.fromRoute) && ['OrderHelp'].indexOf(this.fromRoute.name) !== -1&&this.takeoutId!==0) {  //從訂單明細點修改訂單過來，購物車顯示訂單資料(9)
                    let orderRecord = undefined;
                    this.$sqlite.query('select * from order_record_snapshot_data WHERE takeoutId = ?',[this.takeoutId],(resultSet)=>{
                        if (resultSet.rows.length>0) {
                            orderRecord = JSON.parse(resultSet.rows.item(0).data_json);
                            //購物車畫面上所有變數回到訂單的狀態
                            for (const key in orderRecord) {
                                if (Object.hasOwnProperty.call(orderRecord, key)&&key!=='concatCouponList') {
                                    const element = orderRecord[key];
                                    this[key] = element;

                                }
                            }

                            this.$store.commit(mutationTypes.setOrderStoreInfo, this.storeInfo);
                            this.$store.commit("updateCartPaymentType",this.paymentType);
                            this.$store.commit("updateCartTaxIdNumber",this.taxIdNumber);
                            this.$store.commit("updateShoppingCartMealData",this.shoppingCartMealData);
                            this.$store.commit("updateShoppingCartMealChoiceList",this.shoppingCartMealChoiceList);
                            this.$store.commit("updateCartCouponChoiceList",this.couponChoiceList);
                            this.$store.commit("updateCartcouponChoiceListItem",this.couponChoiceListItem);
                            this.$store.commit("updateCartDiscountLimit",this.discountLimit);
                            //全部餐點選項寫進 order_record_meal_choiceList
                            for (const key in this.shoppingCartMealChoiceList) {
                                if (Object.hasOwnProperty.call(this.shoppingCartMealChoiceList, key)) {
                                    const element = this.shoppingCartMealChoiceList[key];
                                    this.$sqlite.insert('insert into order_record_meal_choiceList (takeoutId,amount,choicePrice,meal_id,combination_json) values (?, ?, ?, ?, ?);',[this.takeoutId,element.amount,element.choicePrice,element.meal_id,element.combination_json],(resultSet)=>{});
                                }
                            }
                            //全部餐點資料寫進 order_record_meal_data
                            for (const key in this.shoppingCartMealData) {
                                if (Object.hasOwnProperty.call(this.shoppingCartMealData, key)) {
                                    const element = this.shoppingCartMealData[key];
                                    this.$sqlite.insert('insert into order_record_meal_data (takeoutId,meal_id,data_json) values(?, ?, ?);',[this.takeoutId,element.meal_id,element.data_json],(resultSet)=>{});
                                }
                            }
                        }
                    });
                }

                if (notNull(this.fromRoute) && ['OrderCouponList','OrderDetail','OrderAddressList','EditOrderList',].indexOf(this.fromRoute.name) === -1) {  // 載入線上點餐相關設定(7)
                    this.$http.fetchWithAuth`GetOnlineOrderAndOtherSetting${{
                    'storeId': this.storeId,
                    }}
                    ${json => {
                    if (json.status) {
                        if (notNull(json.data)) {
                            let data = json.data;
                            this.deliveryConditionList = data.deliveryConditionList.filter(item=>item.status);
                            this.maxDistance = data.maxDistance;

                            this.remarkPrompt = data.storeInfo.remarkPrompt;
                            this.isRequiredRemark = data.storeInfo.isRequiredRemark;
                            this.isRequestTableware = data.storeInfo.isRequestTableware;
                            this.checkoutInstructions = data.storeInfo.checkoutInstructions;
                            this.isCanCoupon = data.storeInfo.isCanCoupon;
                            this.isCanTaxIdNumber = data.storeInfo.isCanTaxIdNumber;
                            this.mealPreparationTimeType = data.storeInfo.mealPreparationTimeType;
                            this.orderDeliveryList = data.storeInfo.orderDeliveryList;

                            if(this.$store.state.orderStoreInfo.reservationTime === '' ||
                                (!this.isCordova && typeof this.$store.state.orderStoreInfo.reservationTime === "undefined")) {
                              this.storeInfo = json.data.storeInfo;
                              this.$store.commit(mutationTypes.setOrderStoreInfo, this.storeInfo);
                            } else {
                              this.storeInfo = this.$store.state.orderStoreInfo
                            }

                            this.radioTake = this.$store.state.orderStoreInfo.reservationType || this.radioTake
                            this.$store.commit("updateCartRadioTake", this.radioTake);
                            if (this.takeoutId === 0) {
                              this.userAddressData = json.data.userAddressData;
                            }
                            this.paymentMethod = json.data.paymentMethod;
                            this.$store.commit("updateCartOnlineOrderAndOtherSetting",{
                                deliveryConditionList:this.deliveryConditionList,
                                maxDistance:this.maxDistance,
                                remarkPrompt:this.remarkPrompt,
                                isRequiredRemark:this.isRequiredRemark,
                                isRequestTableware:this.isRequestTableware,
                                checkoutInstructions:this.checkoutInstructions,
                                isCanCoupon:this.isCanCoupon,
                                isCanTaxIdNumber:this.isCanTaxIdNumber,
                                mealPreparationTimeType:this.mealPreparationTimeType,
                                orderDeliveryList:this.orderDeliveryList,
                                storeInfo:this.storeInfo,
                                userAddressData:this.userAddressData,
                                paymentMethod:this.paymentMethod,
                            });
                            this.$store.commit("updateCartUserAddressId", this.userAddressData.id)
                            this.userAddressId = this.userAddressData.id
                        }
                    }
                    }}`;
                }
            },
            // 確定刪除餐點按鈕(按鈕會觸發彈框關閉事件)
            delMeal(){
                //在準備刪除的陣列新增確定刪除的餐點
                this.delMealArr.push(this.shoppingCartMealChoiceList[this.delMealIndex]);
                //頁面上的陣列移掉刪除的餐點
                this.shoppingCartMealChoiceList.splice(this.delMealIndex,1);
                this.$store.commit("updateShoppingCartMealChoiceList",this.shoppingCartMealChoiceList);

                this.delMealIndex = null;
                this.showDelDialog = false;

                this.$store.commit("trigerCartDelDialog",{
                  showDelDialog: false,
                  delMealIndex: null,
                  delMealOldVal: null,
                });
            },
            // 刪除餐點防呆彈框關閉事件
            closeDelDialog(){
                if (notNull(this.delMealOldVal)&&notNull(this.$refs['MealList'].$refs['setMeal'][this.delMealIndex])) {
                    this.$refs['MealList'].$refs['setMeal'][this.delMealIndex].setNum(this.delMealOldVal);
                }
                this.delMealIndex = null;
                this.delMealOldVal = null;
            },
            // 送出訂單回傳的彈框關閉事件
            closeSuccessDialog(){
                //因為調整成購物車回點餐列表也是快照，所以新增完訂單要清快照
                this.$store.commit("savePageSnapshot",{routeName: "OrderList", variable: undefined});
                if(this.orderSn === '') {
                  this.$router.replace({name: 'OrderRecordList'});
                } else {
                  this.$router.replace({name: 'OrderRecordDetail',query:{orderSn: this.orderSn}});
                }
            },
            updateShopCart(){
                //判斷是否為修改訂單
                let mealChoiceListTable = '';
                let mealDataTable = '';
                if (this.takeoutId !== 0) {
                    mealChoiceListTable = 'order_record_meal_choiceList';
                    mealDataTable = 'order_record_meal_data';
                } else {
                    mealChoiceListTable = 'shopping_cart_meal_choiceList';
                    mealDataTable = 'shopping_cart_meal_data';
                }
                //更新餐點選項
                let hasMealIdArray = [];
                for (const key in this.shoppingCartMealChoiceList) {
                    if (Object.hasOwnProperty.call(this.shoppingCartMealChoiceList, key)) {
                        const element = this.shoppingCartMealChoiceList[key];
                        if(this.isCordova) {
                          this.$sqlite.update('UPDATE '+mealChoiceListTable+' SET amount=?, choicePrice=?, text_note=? WHERE id = ?',[element.amount, element.choicePrice, element.text_note, element.id],function(resultSet){});
                        } else {
                          //桌邊服務目前沒有建立記錄表，所以一律先跑原表
                          this.$sqlite.updateDB("shopping_cart_meal_choiceList", [element],function(resultSet){});
                          hasMealIdArray.push(element.meal_id)
                        }
                    }
                }
                //刪除餐點選項
                this.delMealArr.forEach((element, index) => {
                  if(this.isCordova) {
                    this.$sqlite.delete('delete from ' + mealChoiceListTable + '  where id = ?', [element.id], (resultSet) => {
                      delete this.delMeal[index];
                    });
                  } else {
                     //桌邊服務目前沒有建立記錄表，所以一律先跑原表
                     this.$sqlite.deleteDB("shopping_cart_meal_choiceList", "", element.id, (resultSet) => {
                       delete this.delMeal[index];
                     })
                    }
                });

                //刪除沒有被放在選項裏的餐點資料
                if(this.isCordova) {
                  this.$sqlite.delete('delete from ' + mealDataTable + ' where not exists (select id from ' + mealChoiceListTable + ' where ' + mealChoiceListTable + '.meal_id=' + mealDataTable + '.meal_id)', [], function (resultSet) {});
                } else {
                  //桌邊服務目前沒有建立記錄表，所以一律先跑原表
                  this.$sqlite.queryDB("shopping_cart_meal_data", "", "", (resultSet) => {
                    let length = resultSet.length;
                    if(length > 0) {
                      for(let i = 0;i < length;i++) {
                        let item = resultSet[i];
                        if(hasMealIdArray.indexOf(item.meal_id) === -1) {
                          this.$sqlite.deleteDB("shopping_cart_meal_data", "meal_id", item.meal_id, (resultSet) => {})
                        }
                      }
                    }
                  })
                }
            },
            // api格式餐點資料
            mealApiFormat(){
                let apiFormatData = [];
                // 餐點選項
                this.shoppingCartMealChoiceList.forEach((item)=>{
                    let mealDataDetail = null;
                    // 找到此選項所選餐點資訊
                    for (const key in this.shoppingCartMealData) {
                        if (Object.hasOwnProperty.call(this.shoppingCartMealData, key)) {
                            const element = this.shoppingCartMealData[key];
                            if (element.meal_id===item.meal_id) {
                                mealDataDetail = element;
                            }
                        }
                    }
                    // 找的到就繼續
                    if (notNull(mealDataDetail)) {
                        // 選項組合
                        let combination = JSON.parse(item.combination_json);
                        // 餐點詳細資料
                        let detailList = JSON.parse(mealDataDetail.data_json);
                        let temp = {
                            id:detailList.posFoodInfo.id,
                            name:detailList.posFoodInfo.onlineOrderName,
                            orderType:detailList.posFoodInfo.orderType,
                            amount:item.amount,
                            choicePrice:item.choicePrice,
                            textNote:item.text_note,
                            posFoodSubmealList:function(){
                                if (detailList.posFoodInfo.orderType===2) {
                                    let tempSubMealList = [];
                                    combination[0].forEach((combinationItem,combinationIndex)=>{
                                        let mealSet = detailList.dataList[combinationIndex];
                                        let subMealInfo = mealSet.posFoodSubmealList[combinationItem];
                                        let tempSubMeal = {
                                            mealSetId: mealSet.id,
                                            mealSetName: mealSet.name,
                                            mealSetQuantity: mealSet.quantity,
                                            posFoodId: subMealInfo.posFoodId,
                                            posFoodName: subMealInfo.posFoodName,
                                            increasePrice: subMealInfo.increasePrice,
                                            posFoodSubmealId: subMealInfo.posFoodSubmealId,
                                            // 餐點註記選項
                                            PosGoodsNoteDetailList: function () {
                                                let tempNoteList = [];
                                                combination[1][combinationIndex][combinationItem].forEach((combinationNoteItem,combinationNoteIndex)=>{
                                                    // 過濾掉沒有選的
                                                    if (combinationNoteItem.length>0) {
                                                        let tempNote = subMealInfo.posFoodGoodsNoteList[combinationNoteIndex];
                                                        // 過濾掉沒有選的
                                                        tempNote.posGoodsNoteDetailList = tempNote.posGoodsNoteDetailList.filter((tempNoteItem,tempNoteIndex)=>{return combinationNoteItem.indexOf(tempNoteIndex)!==-1;});
                                                        tempNoteList.push(...tempNote.posGoodsNoteDetailList);
                                                    }
                                                });
                                                return tempNoteList;
                                            }(),
                                        };
                                        tempSubMealList.push(tempSubMeal);
                                    });
                                    return tempSubMealList;
                                }
                            }(),
                            PosGoodsNoteDetailList:function(){
                                if (detailList.posFoodInfo.orderType===1) {
                                    let tempNoteList = [];
                                    combination.forEach((combinationItem,combinationIndex)=>{
                                        // 過濾掉沒有選的
                                        if (combinationItem.length>0) {
                                            let tempNote = detailList.dataList[combinationIndex];
                                            // 過濾掉沒有選的
                                            tempNote.posGoodsNoteDetailList = tempNote.posGoodsNoteDetailList.filter((tempNoteItem,tempNoteIndex)=>{return combinationItem.indexOf(tempNoteIndex)!==-1;});
                                            tempNoteList.push(...tempNote.posGoodsNoteDetailList);
                                        }
                                    });
                                    return tempNoteList;
                                }
                            }(),
                            PosFoodInfo:{
                                id: detailList.posFoodInfo.id,
                                onlineOrderName: detailList.posFoodInfo.onlineOrderName,
                                orderType: detailList.posFoodInfo.orderType,
                                amount: item.amount,
                                choicePrice: item.choicePrice,
                                textNote: item.text_note,
                            },
                        };
                        apiFormatData.push(temp);
                    }
                });
                return apiFormatData;
            },
            changePaymentType(val) {
                this.$store.commit("updateCartPaymentType",val);
                this.paymentType = val;
                this.showBank = false
            },
            creditCardCompleted(){
                /*信用卡結帳後，轉移至訂單完成頁 */
                let self = this;

                //清空購物車
                self.$public.clearCart(this);
                self.$store.commit("clearCartPageSnapshot");
                
                self.$http.fetchWithAuth`PayByAgreedCreditCard${{
                    "orderTmpId": self.orderTmpId,
                }}
                ${json => {
                    if (json.status) {
                        self.orderSn = json.data
                        
                        self.digitalPayModel = false;
                        self.digitalPayShowStr = true;
                        self.completeModel = true;

                        //self.$router.replace({name:'OrderRecordDetail', query:{orderSn:json.data}});
                    } else {
                        self.digitalPayModel = true;
                        self.digitalPayShowStr = false;
                        self.completeModel = false;

                        self.digitalPayErrStr = '抱歉：' + json.message;
                        self.$public.showNotify('抱歉：' + json.message, json.status);
                    }
                }}`;

                /* 高鉅 - 移除
                let inAppBrowserRef = window.cordova.InAppBrowser.open(
                    self.openUrl, '_blank', 'toolbar=no,location=no'
                )
                */

                // inAppBrowserRef.addEventListener('exit', function(){
                //     //關閉視窗要做的事
                // });

                /* 高鉅 - 移除
                inAppBrowserRef.addEventListener('loadstop', function(event) {
                    if(event.url.includes('payment')){
                        inAppBrowserRef.close();
                        localStorage.setItem('routerGoEnable', 1);
                        self.$http.fetch`GetOrderSnByTmpId${{
                            "orderTmpId": self.orderTmpId,
                        }}
                        ${json => {
                            self.$router.push({name:'OrderRecordDetail', query:{orderSn:json.data}});
                        }}`;
                    }
                });
                */
            },
            creditCardOpenWindow() {
                /*開啓藍薪信用卡畫面 */
                let self = this;

                let dataToken = self.$crypto.objectToURIComponent({
                    "orderTmpId":  self.orderTmpId,
                    "paymentType": 2,
                    "itemDesc": "桌邊點餐"
                    })
                let urlSet = process.env.VUE_APP_UTIL_API_HOST + "/newebpayPost" + "?token=" + dataToken + "&pay=2";

                let browser = this.$public.isCordova()? window.cordova.InAppBrowser: window;
                
                if (this.$public.isCordova()) {
                    let appInBrowser = browser.open(urlSet, '_blank', 'location=no,toolbarcolor=#ffffff,closebuttoncaption=關閉,closebuttoncolor=#000000');

                    appInBrowser.addEventListener('exit', function(event){
                        // 監聽 '視窗關閉' 事件
                        self.creditCardCloseWindow()
                    })
                } else {
                    let appInBrowser = browser.open(urlSet, '_blank', 'toolbar=no,location=no');

                    // 迴圈確認 '視窗是否關閉'
                    var loop = setInterval(function() { 
                        if(appInBrowser.closed) {
                            clearInterval(loop)
                            self.creditCardCloseWindow()
                        } 
                    }, 1000);
                }
            },
            creditCardCloseWindow() {
                /* 關閉藍薪信用卡畫面 */
                let self = this;

                self.$http.fetchWithAuth`PayByCreditCardCompleted${{
                    "orderTmpId": self.orderTmpId,
                }}
                ${json => {
                    if (json.status) {
                        //清空購物車
                        self.$public.clearCart(this);
                        self.$store.commit("clearCartPageSnapshot");

                        if(self.isCordova || self.isOnlineOrder) {
                            self.orderSn = json.data.orderSn
                        
                            self.digitalPayModel = false;
                            self.digitalPayShowStr = true;
                            self.completeModel = true;
                        } else {
                            self.dataObject = {
                                ...self.dataObject, ...{
                                    "takeoutId": json.data.takeoutId,
                                }
                            }
                            sessionStorage.setItem(constDefine.StorageNameType.TableOrderDataToken, self.$crypto.objectToURIComponent(self.dataObject));

                            let dataToken = self.$crypto.objectToURIComponent({
                                "takeoutId":        json.data.takeoutId,
                                "storeId":          self.dataObject.storeId,
                                "prototypeObject":  {...self.dataObject.prototypeObject},
                            })

                            if (self.dataObject.routePath) {
                                self.$router.push({path: self.dataObject.routePath + "/complete", query: {dataToken: dataToken}});
                            } else {
                                self.$router.push({name: 'TableOrderComplete', query: {dataToken: dataToken}});
                            }
                        }
                        
                    } else {
                        self.digitalPayModel = true;
                        self.digitalPayShowStr = false;
                        self.completeModel = false;

                        self.digitalPayErrStr = '抱歉：' + json.message;
                        self.$public.showNotify('抱歉：' + json.message, json.status);
                    }
                }}`;
            },
            checkCouponChoiceListItem() {
              let isUpdate = false

              Object.keys(this.checkCumulativeTypeList).forEach(index => {
                let item = this.checkCumulativeTypeList[index]
                item.isUpdate = false
                if(item.posUseCumulativeType === 1 && item.posUseAmount > this.originTotalPrice) {
                  isUpdate = true
                  item.isUpdate = true
                  item.canUseAmount = 0
                  item.totalChoiceAmount = 0
                } else if(item.posUseCumulativeType === 2) {
                  let count = Math.floor(this.originTotalPrice / item.posUseAmount)
                  let canUseAmount = count * item.posUsePieces
                  if(item.totalChoiceAmount > canUseAmount) {
                    isUpdate = true
                    item.isUpdate = true
                    item.canUseAmount = canUseAmount
                    item.totalChoiceAmount = canUseAmount
                  }
                }
              })

              if(isUpdate) {
                let tmpCouponChoiceListItem = JSON.parse(JSON.stringify(this.couponChoiceListItem))
                let newtCouponChoiceListItem = []

                tmpCouponChoiceListItem.forEach(item => {
                  let couponListId = "couponListId_"+item.couponListId
                  let checkItem = this.checkCumulativeTypeList[couponListId]
                  if(checkItem.isUpdate) {

                    if(checkItem.canUseAmount === 0) {
                      return false
                    } else {
                      if(checkItem.canUseAmount >= item.choiceAmount) {
                        checkItem.canUseAmount -= item.choiceAmount
                        newtCouponChoiceListItem.push(item)
                      } else {
                        item.choiceAmount = checkItem.canUseAmount
                        checkItem.canUseAmount = 0
                        newtCouponChoiceListItem.push(item)
                      }
                    }
                  } else {
                    newtCouponChoiceListItem.push(item)
                  }
                })

                this.couponChoiceListItem = newtCouponChoiceListItem
                this.$store.commit("updateCartcouponChoiceListItem", this.couponChoiceListItem);
              }
            },
        },
        watch:{
            shoppingCartMealChoiceList:{
                handler:function (newVal,oldVal) {
                    let sumPrice = 0;
                    let sumAmount = 0;
                    for (const key in newVal) {
                        if (Object.hasOwnProperty.call(newVal, key)) {
                            const element = newVal[key];
                            sumPrice += element.choicePrice;
                            sumAmount += element.amount;
                        }
                    }
                    this.$store.commit("updateCartOriginTotalPrice", sumPrice)
                    this.originTotalPrice = sumPrice;
                    this.sumAmount = sumAmount;
                },
                deep : true ,
                immediate : true,
            },
            sumAmount:{
                handler:function (newVal,oldVal) {
                    if (newVal===0&&oldVal!==undefined&&oldVal!==0) {
                        // 購物車無商品自動導回上一頁
                        this.$router.go(-1);
                    }
                },
                deep: true,
                immediate: true,
            },
            orderTmpId:{
                handler:function (newVal, oldVal) {
                    if (this.paymentType == 2) {
                        this.creditCardOpenWindow();
                    } else if (this.paymentType == 8) {
                        this.creditCardCompleted();
                    }
                },
                deep : true ,
                immediate : false,
            },
            showBank:{
                handler:function (newVal,oldVal) {
                    this.$store.commit("updateCartShowBank",newVal);
                },
                deep:false,
                immediate:false,
            },
            '$store.state.clearCartData.isClear': {
                //電子錢包付款完成後，轉至交易成功頁
                handler:function (newVal, oldVal) {
                    //轉頁前，清空購物車
                    this.$public.clearCart(this);
                    this.$store.commit("clearCartPageSnapshot");

                    if (this.$store.state.clearCartData.action == "OCP電子錢包") {
                        this.orderSnOCP = this.$store.state.clearCartData.orderSn;
                        this.$store.commit('setClearCartData',{
                            action: '',
                            isClear: false,
                            orderSn: '',
                        });
                        this.$router.replace({name: 'OrderRecordDetail',query:{orderSn: this.orderSnOCP}});
                    }
                },
                deep : true ,
                immediate : false,
            },
        },
        computed:{
            //訂單即將送出的最終金額
            finalPrice:function () {
                let temp = this.originTotalPrice;

                this.checkCouponChoiceListItem()

                if (this.calculatorCouponPrice !== null) {
                    temp = this.calculatorCouponPrice(temp, this.couponChoiceListItem);

                    let discountLimitPrice = 0;

                    if (notNull(this.discountLimit)) {
                      //自取
                      if (Number(this.radioTake) === constDefine.RADIO_TAKE.takeout) {
                          //金額限制
                          if (this.discountLimit.takeoutDiscountLimitType === 1) {
                            discountLimitPrice = this.discountLimit.takeoutDiscountLimitPrice;
                          //比例限制
                          } else if(this.discountLimit.takeoutDiscountLimitType === 2) {
                            discountLimitPrice = Math.round(this.originTotalPrice * (this.discountLimit.takeoutDiscountLimitPercentage / 100));
                          }
                        //外帶
                      } else if (Number(this.radioTake)===constDefine.RADIO_TAKE.delivery) {
                        //金額限制
                        if (this.discountLimit.deliveryDiscountLimitType === 1) {
                          discountLimitPrice = this.discountLimit.deliveryDiscountLimitPrice;
                        //比例限制
                        } else if(this.discountLimit.deliveryDiscountLimitType === 2) {
                          discountLimitPrice = Math.round(this.originTotalPrice * (this.discountLimit.deliveryDiscountLimitPercentage / 100));
                        }
                      }
                    }

                    if (discountLimitPrice !== 0 && (this.originTotalPrice - temp > discountLimitPrice)) {
                        temp = this.originTotalPrice - discountLimitPrice;
                    }
                }

                if (this.radioTake==='2') {
                    temp += this.deliveryFee;
                }

                this.$store.commit("updateCartFinalPrice",temp)
                return temp;
            },
            //送出訂單的執行事件
            submitToDo:function () {
                let self = this;

                return function () {
                    if (self.paymentType === "") {
                        self.$el.style.setProperty('--borderColor333', `${"#EA4E5E"}`)
                        document.getElementById('PayType').scrollIntoView({behavior: "smooth", block: "center", inline: "nearest"});
                        self.$public.showNotify("請選擇付款方式", false);
                        return false
                    }   else if (!self.storeInfo.orderStatus && self.storeInfo.reservationTime === "") {
                        self.$public.showNotify("目前僅接收預約訂單", false);
                        self.$store.commit(mutationTypes.setOpenPickTime, true);
                        return false
                    }

                    self.isLoading = true;

                    this.$http.fetchWithAuthEncrypt`CreateOrder2${{
                        storeId: self.storeInfo.id,
                        mealChoiceList: self.shoppingCartMealChoiceList,
                        mealData: self.shoppingCartMealData,
                        remark: self.remark,
                        sumAmount: self.sumAmount,
                        couponChoiceList: self.couponChoiceListItem,
                        taxIdNumber: self.taxIdNumber,
                        userAddressId: self.userAddressId,
                        radioTake: parseInt(self.radioTake),
                        isRequiredTableware: self.isRequiredTableware,
                        totalPrice: self.finalPrice,
                        takeoutId: self.takeoutId,
                        paymentType: parseInt(self.paymentType),
                        reservationTime: self.storeInfo.reservationTime,
                        mealApiFormat:self.mealApiFormat(),
                        isCordova: self.isCordova,
                        phone: self.phone,
                        tableId: self.tableId,
                        adult: self.adult,
                        child: self.child,
                    }}
                    ${json => {
                        self.$public.showNotify(json.message, json.status);
                        self.isLoading = false
                        if (json.status) {
                            if (self.paymentType === '1') {
                                //記住送成功前的購物車頁面資訊
                                let recordCount = 0;
                                let dataJson = {};
                                this.$public.clearCart(self);
                                this.$store.commit("clearCartPageSnapshot");
                                if (json.data.takeoutId!==0) {
                                    self.takeoutId = json.data.takeoutId;
                                }
                                for (const key in self._data) {
                                    if (Object.hasOwnProperty.call(self._data, key)&&key!=="fromRoute") {
                                      const element = self._data[key];
                                      dataJson[key] = element;
                                    }
                                }
                                dataJson = JSON.stringify(dataJson);

                                if(self.isCordova || self.isOnlineOrder) {
                                  self.completeModel = true
                                  self.orderSn = json.data.orderSn
                                  self.completeOrderData.userName = json.data.userName
                                  self.completeOrderData.radioTake = json.data.radioTake
                                  self.completeOrderData.totalPrice = json.data.totalPrice
                                  self.completeOrderData.day = json.data.day
                                  self.completeOrderData.time = json.data.time
                                  self.completeOrderData.week = json.data.week
                                  //確認此訂單是否有存在暫存
                                  self.$sqlite.query('select * from order_record_snapshot_data WHERE takeoutId = ?', [Number(self.takeoutId)], (resultSet) => {
                                        recordCount = resultSet.rows.length;
                                      //有的話就更新，沒有的話就新增
                                      if (recordCount > 0) {
                                          self.$sqlite.update('UPDATE order_record_snapshot_data SET data_json=? WHERE takeoutId = ?', [dataJson, self.takeoutId], (resultSet) => {
                                          });
                                      } else {
                                          self.$sqlite.insert('insert into order_record_snapshot_data (orderId,orderSn,takeoutId, data_json) values (?, ?, ?, ?);', [json.data.orderId, json.data.orderSn, json.data.takeoutId, dataJson], (resultSet) => {
                                          });
                                      }
                                  });
                                } else {
                                  self.dataObject = {
                                    ...self.dataObject, ...{
                                      "takeoutId": json.data.takeoutId,
                                    }
                                  }
                                  sessionStorage.setItem(constDefine.StorageNameType.TableOrderDataToken, self.$crypto.objectToURIComponent(self.dataObject));

                                  let dataToken = self.$crypto.objectToURIComponent({
                                    "takeoutId":        json.data.takeoutId,
                                    "storeId":          self.dataObject.storeId,
                                    "prototypeObject":  {...self.dataObject.prototypeObject},
                                  })

                                  if (self.dataObject.routePath) {
                                    self.$router.push({path: self.dataObject.routePath + "/complete", query: {dataToken: dataToken}});
                                  } else {
                                    self.$router.push({name: 'TableOrderComplete', query: {dataToken: dataToken}});
                                  }
                                }
                            } else if (self.paymentType === '2' || self.paymentType === '8') {
                                self.digitalPayModel = true
                                self.digitalPayShowStr = true;
                                self.completeModel = false;

                                this.$public.clearCart(self);
                                //信用卡使用的暫存訂單編號
                                self.orderTmpId += json.data.orderTmpId;
                                //付款完成彈出視窗內容
                                self.completeOrderData.userName = json.data.userName
                                self.completeOrderData.radioTake = json.data.radioTake
                                self.completeOrderData.totalPrice = json.data.totalPrice
                                self.completeOrderData.day = json.data.day
                                self.completeOrderData.time = json.data.time
                                self.completeOrderData.week = json.data.week
                            }
                        }
                    }}`;
                }
            },
        },
        beforeRouteEnter: function(to, from, next) {
            next(
                vm=>{
                    // vm就是出去的this，先把來的路由存起來，出去就抓不到了？
                    vm.fromRoute = from;
                }
            );
        },
        beforeRouteLeave: function(to, from, next) {
            this.updateShopCart();
            // 如果是'OrderDetail','OrderCouponList'的話就記住快照
            if (['OrderDetail','OrderCouponList','OrderAddressList','EditOrderList'].indexOf(to.name) !== -1) {
                let variable = {};
                for (const key in this._data) {
                    if (Object.hasOwnProperty.call(this._data, key)) {
                        const element = this._data[key];
                        variable[key] = element;
                    }
                }
                variable.scrollTop = document.querySelector('.main').scrollTop;
                variable.scrollHeight = document.querySelector('.main').scrollHeight;
                this.$store.commit("saveAllCartPageSnapshot",variable);
            }
            if (['OrderHelp'].indexOf(to.name) !== -1) {
                if(this.isCordova) {
                  this.$sqlite.delete('delete from order_record_meal_choiceList',[],function(resultSet){});
                  this.$sqlite.delete('delete from order_record_meal_data',[],function(resultSet){});
                }
            }
            next();
        },
    }
</script>
<style lang="less" scoped src="../../less/order.less"/>
