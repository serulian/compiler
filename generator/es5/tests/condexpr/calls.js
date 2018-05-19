$module('calls', function () {
  var $static = this;
  this.$class('165552b1', 'SomeError', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Message = $t.property(function () {
      var $this = this;
      return $t.fastbox('huh?', $g.________testlib.basictypes.String);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Message|3|b2b53db7": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.getValue = function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  };
  $static.failValue = function () {
    throw $g.calls.SomeError.new();
  };
  $static.getIntValue = function () {
    return $t.fastbox(45, $g.________testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    return $g.calls.getIntValue().$wrapped == 2 ? $g.calls.failValue() : $g.calls.getValue();
  };
});
