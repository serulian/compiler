$module('with', function () {
  var $static = this;
  this.$class('20c15ffe', 'SomeReleasable', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Release = function () {
      var $this = this;
      $g.with.someBool = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      return;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Release|2|29dc432d<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $temp0;
    var $resources = $t.resourcehandler();
    $t.fastbox(123, $g.____testlib.basictypes.Integer);
    $temp0 = $g.with.SomeReleasable.new();
    $resources.pushr($temp0, '$temp0');
    $t.fastbox(456, $g.____testlib.basictypes.Integer);
    $resources.popr('$temp0');
    $t.fastbox(789, $g.____testlib.basictypes.Integer);
    var $pat = $g.with.someBool;
    $resources.popall();
    return $pat;
  };
  this.$init(function () {
    return $promise.new(function (resolve) {
      $static.someBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      resolve();
    });
  }, '19a53bdb', []);
});
