$module('this', function () {
  var $static = this;
  this.$class('93ce2e23', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.DoSomething = function () {
      var $this = this;
      return;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "DoSomething|2|89b8f38e<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

});
