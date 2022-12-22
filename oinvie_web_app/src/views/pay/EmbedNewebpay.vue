<template>
    <div>
        <navigation-bar Back/>
        <div class="main container">
            <!-- 這個是 A組 -->
            <template v-if="!haveToken">
                <div class="no-data">
                    <img src="../../assets/images/empty.svg">
                    <p>目前尚無綁定信用卡</p>
                </div>
                <div class="footer-button mt-0 mb-5">
                    <el-button type="primary" class="w-50" round @click="click">+ 綁定信用卡</el-button>
                </div>
            </template>

            <!-- 這個是 B組 -->
            <template v-else>
                <small>信用卡卡號</small>
                <div class="dark">{{creditCardNumber}}</div>
                <hr>
                <small>授權有效日期</small>
                <div class="dark">{{tokenLife}}</div>
                <hr>
                <div class="footer-button mb-5">
                    <!-- <el-button class="w-50" round @click="click = true">刪除綁定</el-button> -->
                    <el-button type="primary" class="w-50" round @click="click">重新綁定</el-button>
                </div>
            </template>

            <!-- A、B 不會同時出現 -->
        </div>

        <loading
            :active.sync=isLoading
            :can-cancel="false"
            :is-full-page="true">
        </loading>
    </div>
</template>

<script>

    export default {
        name: "EmbedNewebpay",
        data() {
            return {
                isLoading:false,

                iframeSrc: '',
                haveToken: false,

                tokenLife: '',
                creditCardNumber: '',
            }
        },
        created() {
            this.setupData();
        },
        mounted() {
            
        },
        methods: {
            setupData() {
                this.$http.fetchWithAuth`UserPayTokenByType${{
                    'tokenType': 1,
                }}
                ${json => {
                    if (json.status) {
                        this.haveToken = true;
                        this.tokenLife = json.data.tokenLife;
                        this.creditCardNumber = json.data.creditCardNumber;
                    } else {
                        this.haveToken = false;
                        this.tokenLife = "1900/01/01";
                        this.creditCardNumber = "000000******0000";
                    }
                }}`;
            },
            click() {
                let self = this;
                self.iframeSrc = process.env.VUE_APP_UTIL_API_HOST + "/newebpayPost" + "?token=" + encodeURIComponent(window.localStorage.getItem('token'));
                
                //參考文件 https://cordova.apache.org/docs/en/10.x/reference/cordova-plugin-inappbrowser/
                let inAppBrowserRef = window.cordova.InAppBrowser.open(
                    self.iframeSrc, '_blank', 'location=no,toolbarcolor=#ffffff,closebuttoncaption=關閉,closebuttoncolor=#000000'
                )
                self.isLoading = true;

                inAppBrowserRef.addEventListener('exit', function(event){
                    // 監聽 '視窗關閉' 事件
                    self.setupData();
                    self.isLoading = false;
                });
            },
        }

    }
</script>

<style lang="less" scoped src="../../less/pay.less"/>
