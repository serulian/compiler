$module('basic', function () {
  var $static = this;
  this.$class('7dbf4efd', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeInt = $t.fastbox(2, $g.____testlib.basictypes.Integer);
      instance.AnotherBool = $g.basic.CoolFunction();
      return instance;
    };
    $instance.AnotherFunction = function () {
      var $this = this;
      return;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherFunction|2|29dc432d<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.CoolFunction = function () {
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.basic.SomeClass.new().AnotherBool;
  };
});
